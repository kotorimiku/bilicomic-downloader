package bilicomicdownloader

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MetaData struct {
	Title       string
	Creator     *string
	Publisher   *string
	Description *string
	Series      *string
	Subject     []string
	Language    *string
	Index       *int
	Identifier  *string
}

func NewMetaData(
	title string,
	creator *string,
	publisher *string,
	description *string,
	series *string,
	subject []string,
	language *string,
	index *int,
	identifier *string,
) MetaData {
	title = escapeEpubText(title)
	if creator != nil {
		*creator = escapeEpubText(*creator)
	}
	if publisher != nil {
		*publisher = escapeEpubText(*publisher)
	}
	if description != nil {
		*description = escapeEpubText(*description)
	}
	if series != nil {
		*series = escapeEpubText(*series)
	}
	if language != nil {
		*language = escapeEpubText(*language)
	}
	if identifier != nil {
		*identifier = escapeEpubText(*identifier)
	}
	return MetaData{
		Title:       title,
		Creator:     creator,
		Publisher:   publisher,
		Description: description,
		Series:      series,
		Subject:     subject,
		Language:    language,
		Index:       index,
		Identifier:  identifier,
	}
}

type EpubBuilder struct {
	metadata    MetaData
	text        []string
	chapterList []string
	imgDataList [][]byte
	extList     []string
	imgPathList []string
	addCatalog  bool
}

func NewEpubBuilder(
	metadata MetaData,
	text []string,
	chapterList []string,
	imgDataList [][]byte,
	extList []string,
	addCatalog bool,
) EpubBuilder {
	return EpubBuilder{
		metadata:    metadata,
		text:        text,
		chapterList: chapterList,
		imgDataList: imgDataList,
		extList:     extList,
		addCatalog:  addCatalog,
	}
}

func (eb *EpubBuilder) BuildEpub() map[string][]byte {
	epub := make(map[string][]byte)
	// mimetype needs to be the first file
	epub["mimetype"] = []byte("application/epub+zip")
	epub["META-INF/container.xml"] = []byte(eb.buildContainer())
	epub["OEBPS/content.opf"] = []byte(eb.buildOpf())
	epub["OEBPS/toc.ncx"] = []byte(eb.buildNcx())
	epub["OEBPS/Text/cover.xhtml"] = []byte(eb.buildCoverXhtml())
	for i := 0; i < len(eb.text); i++ {
		epub[fmt.Sprintf("OEBPS/Text/%s.xhtml", eb.numFill(i+1))] = []byte(eb.buildXhtml(eb.chapterList[i], eb.text[i]))
	}
	epub["OEBPS/Text/nav.xhtml"] = []byte(eb.buildNavXhtml())
	for i := 0; i < len(eb.extList); i++ {
		ext := strings.Split(eb.extList[i], ".")[len(strings.Split(eb.extList[i], "."))-1]
		epub[fmt.Sprintf("OEBPS/Images/%s.%s", eb.numFill(i), ext)] = eb.imgDataList[i]
	}
	if eb.addCatalog {
		filePath, fileData := eb.buildSgcNavCss()
		eb.addFile(&epub, filePath, fileData)
	}
	return epub
}

func (eb *EpubBuilder) BuildComic(path string, imgPath string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	zipFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	mimetype := &zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store,
	}

	// Write mimetype first
	mimetypeFile, err := zipWriter.CreateHeader(mimetype)
	if err != nil {
		return err
	}
	mimetypeFile.Write([]byte("application/epub+zip"))

	eb.text = make([]string, 0, 300)
	eb.imgPathList = make([]string, 0, 300)
	err = filepath.Walk(imgPath, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if info.IsDir() {
			return nil
		}

		// 创建一个文件在压缩包中
		relPath, err := filepath.Rel(imgPath, file)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(relPath)

		imgPath := filepath.ToSlash(filepath.Join("OEBPS/Images", relPath))

		writer, err := zipWriter.Create(imgPath)
		if err != nil {
			return err
		}

		fileReader, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		_, err = io.Copy(writer, fileReader)

		if err != nil {
			return err
		}

		eb.text = append(eb.text, eb.BuildComicTag(fmt.Sprintf("../Images/%s", relPath)))

		eb.imgPathList = append(eb.imgPathList, relPath)
		return err
	})

	epub := make(map[string][]byte)
	epub["META-INF/container.xml"] = []byte(eb.buildContainer())
	epub["OEBPS/content.opf"] = []byte(eb.buildOpf())
	epub["OEBPS/Text/nav.xhtml"] = []byte(eb.buildNavXhtml())
	for i := 0; i < len(eb.text); i++ {
		epub[fmt.Sprintf("OEBPS/Text/%s.xhtml", eb.numFill(i+1))] = []byte(eb.buildXhtml("", eb.text[i]))
	}

	for fileName, fileData := range epub {
		if fileName == "mimetype" {
			continue
		}
		file, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}
		file.Write(fileData)
	}

	return err
}

func (eb *EpubBuilder) BuildComicTag(imgPath string) string {
	return fmt.Sprintf(`
<img src="%s" alt="%s"/>`, imgPath, imgPath)
}

func (eb *EpubBuilder) SaveFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	fileMap := eb.BuildEpub()

	zipFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Write mimetype first
	mimetypeFile, err := zipWriter.Create("mimetype")
	if err != nil {
		return err
	}
	mimetypeFile.Write(fileMap["mimetype"])

	for fileName, fileData := range fileMap {
		if fileName == "mimetype" {
			continue
		}
		file, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}
		file.Write(fileData)
	}

	return nil
}

func (eb *EpubBuilder) buildNcx() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN"
 "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
<ncx version="2005-1" xmlns="http://www.daisy.org/z3986/2005/ncx/">
  <head>
    <meta name="dtb:depth" content="1" />
    <meta name="dtb:totalPageCount" content="0" />
    <meta name="dtb:maxPageNumber" content="0" />
  </head>
  <docTitle>
    <text>%s</text>
  </docTitle>
  <navMap>
    %s
  </navMap>
</ncx>`, eb.metadata.Title, eb.getNavXml())
}

func (eb *EpubBuilder) getNavXml() string {
	var navMap []string
	for i := 0; i < len(eb.chapterList); i++ {
		navMap = append(navMap, fmt.Sprintf(`<navPoint id="navPoint-%d" playOrder="%d">
      <navLabel>
        <text>%s</text>
      </navLabel>
      <content src="%s" />
    </navPoint>`, i+1, i+1, eb.chapterList[i], fmt.Sprintf("Text/%s.xhtml", eb.numFill(i+1))))
	}
	return strings.Join(navMap, "\n    ")
}

func (eb *EpubBuilder) buildOpf() string {
	metadata := eb.getMetadataXml()
	manifest := eb.getManifestXml()
	spine := eb.getSpineXml()
	// guide := eb.getGuideXml()
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<package version="3.0" unique-identifier="BookId" xmlns="http://www.idpf.org/2007/opf">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
    %s
  </metadata>
  <manifest>
    %s
  </manifest>
  <spine>
    %s
  </spine>
</package>`, metadata, manifest, spine)
}

func (eb *EpubBuilder) getGuideXml() string {
	if len(eb.extList) == 0 {
		return ""
	}
	return `<reference href="Text/cover.xhtml" title="Cover" type="cover"/>`
}

func (eb *EpubBuilder) getSpineXml() string {
	var spine []string
	if len(eb.extList) > 0 {
		spine = append(spine, `<itemref idref="cover.xhtml"/>`)
	}
	for i := 0; i < len(eb.text); i++ {
		spine = append(spine, fmt.Sprintf(`<itemref idref="x%s.xhtml"/>`, eb.numFill(i+1)))
	}
	return strings.Join(spine, "\n    ")
}

func (eb *EpubBuilder) getManifestXml() string {
	var manifest []string
	if len(eb.extList) > 0 {
		manifest = append(manifest, `<item id="cover.xhtml" href="Text/cover.xhtml" media-type="application/xhtml+xml"/>`)
	}
	// manifest = append(manifest, `<item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>`)

	for i := 0; i < len(eb.text); i++ {
		manifest = append(manifest, fmt.Sprintf(`<item id="x%s.xhtml" href="Text/%s.xhtml" media-type="application/xhtml+xml"/>`, eb.numFill(i+1), eb.numFill(i+1)))
	}

	for i := 0; i < len(eb.imgDataList); i++ {
		ext := strings.Split(eb.extList[i], ".")[len(strings.Split(eb.extList[i], "."))-1]
		manifest = append(manifest, fmt.Sprintf(`<item id="x%s.%s" href="Images/%s.%s" media-type="image/jpeg"/>`, eb.numFill(i), ext, eb.numFill(i), ext))
	}
	for i, imgPath := range eb.imgPathList {
		manifest = append(manifest, fmt.Sprintf(`<item id="x%s" href="Images/%s" media-type="image/jpeg"/>`, fmt.Sprintf("%03d", i), imgPath))
	}
	manifest = append(manifest, `<item id="nav.xhtml" href="Text/nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>`)
	if eb.addCatalog {
		manifest = append(manifest, `<item id="sgc-nav.css" href="Styles/sgc-nav.css" media-type="text/css"/>`)
	}
	return strings.Join(manifest, "\n    ")
}

func (eb *EpubBuilder) getMetadataXml() string {
	var metadata []string
	metadata = append(metadata, fmt.Sprintf(`<dc:title>%s</dc:title>`, eb.metadata.Title))
	if eb.metadata.Creator != nil {
		metadata = append(metadata, fmt.Sprintf(`<dc:creator>%s</dc:creator>`, *eb.metadata.Creator))
	}
	if eb.metadata.Publisher != nil {
		metadata = append(metadata, fmt.Sprintf(`<dc:publisher>%s</dc:publisher>`, *eb.metadata.Publisher))
	}
	if eb.metadata.Description != nil {
		metadata = append(metadata, fmt.Sprintf(`<dc:description>%s</dc:description>`, *eb.metadata.Description))
	}
	if eb.metadata.Language != nil {
		metadata = append(metadata, fmt.Sprintf(`<dc:language>%s</dc:language>`, *eb.metadata.Language))
	} else {
		metadata = append(metadata, `<dc:language>zh</dc:language>`)
	}
	if eb.metadata.Identifier != nil {
		metadata = append(metadata, fmt.Sprintf(`<dc:identifier id="BookId">%s</dc:identifier>`, *eb.metadata.Identifier))
	} else {
		metadata = append(metadata, `<dc:identifier id="BookId">BookId</dc:identifier>`)
	}
	for _, subject := range eb.metadata.Subject {
		metadata = append(metadata, fmt.Sprintf(`<dc:subject>%s</dc:subject>`, subject))
	}
	// metadata = append(metadata, strings.Join(eb.metadata.Subject, "\n\t\t"))
	metadata = append(metadata, fmt.Sprintf(`<meta property="dcterms:modified">%s</meta>`, getTime()))
	if len(eb.extList) > 0 {
		metadata = append(metadata, fmt.Sprintf(`<meta name="cover" content="x000.%s"/>`, strings.TrimPrefix(eb.extList[0], ".")))
	}
	if eb.metadata.Series != nil {
		metadata = append(metadata, fmt.Sprintf(`<meta name="calibre:series" content="%s"/>`, *eb.metadata.Series))
	}
	if eb.metadata.Index != nil {
		metadata = append(metadata, fmt.Sprintf(`<meta name="calibre:series_index" content="%d"/>`, *eb.metadata.Index))
	}
	return strings.Join(metadata, "\n    ")
}

func (eb *EpubBuilder) buildContainer() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml" />
  </rootfiles>
</container>`
}

func (eb *EpubBuilder) buildXhtml(title, body string) string {
	titleTag := ""
	if title != "彩页" && title != "" {
		titleTag = fmt.Sprintf("<h1>%s</h1>\n    ", title)
	} else if title == "" {
		title = "title"
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
  <head>
    <title>%s</title>
    <style type="text/css">p{text-indent:2em;}</style>
  </head>
  <body>
    %s%s
  </body>
</html>`, title, titleTag, body)
}

func (eb *EpubBuilder) numFill(num int) string {
	return fmt.Sprintf("%03d", num)
}

func (eb *EpubBuilder) buildCoverXhtml() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
  <title>Cover</title>
</head>
<body>
  <div style="text-align: center; padding: 0pt; margin: 0pt;">
    <img src="../Images/000%s" alt="cover" />
  </div>
</body>
</html>`, eb.extList[0])
}

func (eb *EpubBuilder) buildNavXhtml() string {
	css := ""
	if eb.addCatalog {
		css = fmt.Sprintf("\n  %s", `<link href="../Styles/sgc-nav.css" rel="stylesheet" type="text/css"/>`)
	}
	var navMap []string
	for i := 0; i < len(eb.chapterList); i++ {
		navMap = append(navMap, fmt.Sprintf("<li><a href=\"%s.xhtml\">%s</a></li>", eb.numFill(i+1), eb.chapterList[i]))
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="en" xml:lang="en">
<head>
  <title>ePub NAV</title>
  <meta charset="utf-8"/>%s
</head>
<body epub:type="frontmatter">
  <nav epub:type="toc" id="toc" role="doc-toc">
    <h1>目录</h1>
    <ol>
      %s
    </ol>
  </nav>
</body>
</html>`, css, strings.Join(navMap, "\n      "))
}

func (eb *EpubBuilder) buildSgcNavCss() (string, []byte) {
	filePath := "OEBPS/Styles/sgc-nav.css"
	return filePath, []byte(`nav#toc {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  padding: 20px;
  background-color: #f8f8f8; /* 浅灰色背景 */
  border-radius: 10px;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1); /* 柔和的阴影 */
}

nav#toc h1 {
  font-size: 24px;
  color: #333;
  text-align: center;
  margin-bottom: 20px;
  font-weight: bold; /* 加粗 */
}

nav#toc ol {
  list-style-type: none;
  padding-left: 0;
}

nav#toc ol li {
  margin-bottom: 10px;
}

nav#toc ol li a {
  text-decoration: none;
  font-size: 18px;
  color: #555;
  padding: 6px;
  display: block;
  transition: background-color 0.3s, color 0.3s;
  border-radius: 5px;
}

nav#toc ol li a:hover {
  background-color: #d9d9d9;
  color: #000;
}
`)
}

func (eb *EpubBuilder) addFile(epub *map[string][]byte, filePath string, fileData []byte) {
	(*epub)[filePath] = fileData
}

func getTime() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func escapeEpubText(input string) string {
	input = strings.ReplaceAll(input, "&", "&amp;")
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	return input
}
