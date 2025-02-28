package bilicomicdownloader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DownloaderSingle struct {
	BookInfo    *BookInfo
	Volume      *Volume
	Index       int
	config      *Config
	Progress    float32
	Fail        bool
	messageSend func(string)
}

func NewDownloaderSingle(bookInfo *BookInfo, volume *Volume, index int, config *Config, messageSend func(string)) *DownloaderSingle {
	return &DownloaderSingle{
		BookInfo:    bookInfo,
		Volume:      volume,
		Index:       index,
		config:      config,
		messageSend: messageSend,
	}
}

func (d *DownloaderSingle) GetImageUrlList(urlChapter string) ([]string, error) {
	text, err := GetRaw(urlChapter)
	if err != nil {
		return nil, err
	}

	html, err := goquery.NewDocumentFromReader(text)
	if err != nil {
		return nil, err
	}

	acontentz := html.Find("#acontentz").First()

	var imageUrls = make([]string, 0, 70)
	acontentz.Find("img").Each(func(i int, s *goquery.Selection) {
		imageUrls = append(imageUrls, s.AttrOr("data-src", ""))
	})

	if len(imageUrls) == 0 {
		return nil, fmt.Errorf("no image found")
	}

	return imageUrls, nil
}

func (d *DownloaderSingle) GetImageUrlListWithRetry(urlChapter string) ([]string, error) {
	for _ = range 10 {
		imgList, err := d.GetImageUrlList(urlChapter)
		if err != nil {
			d.messageSend(fmt.Sprintf("获取图片失败：%v", err.Error()))
			time.Sleep(3 * time.Second)
			continue
		}
		return imgList, nil
	}
	return nil, fmt.Errorf("获取图片失败")
}

func (d *DownloaderSingle) Download(processSend func()) error {
	chapters := d.Volume.Chapters
	folderPath := filepath.Join(d.config.OutputPath, sanitizeFilename(d.BookInfo.Title), sanitizeFilename(d.Volume.Title))
	os.MkdirAll(folderPath, os.ModePerm)
	total := 0

	wg := sync.WaitGroup{}
	maxConcurrency := 1
	sem := make(chan struct{}, maxConcurrency)
	var error error

	for index, chapter := range chapters {
		url := fmt.Sprintf("https://%s%s", d.config.UrlBase, chapter.Url)
		imgList, err := d.GetImageUrlListWithRetry(url)
		if err != nil {
			d.Fail = true
			return err
		}

		total += len(imgList)
		chapterPath := filepath.Join(folderPath, sanitizeFilename(fmt.Sprintf("%d-%s", index+1, chapter.Title)))
		os.MkdirAll(chapterPath, os.ModePerm)

		for i, imgUrl := range imgList {
			sem <- struct{}{}
			wg.Add(1)
			go func(i int, img string) {
				defer func() {
					<-sem
					wg.Done()
				}()
				filePath := filepath.Join(chapterPath, fmt.Sprintf("%03d.%s", i+1, strings.Split(imgUrl, ".")[len(strings.Split(imgUrl, "."))-1]))
				_, err := os.Stat(filePath)
				if os.IsNotExist(err) {
					err := d.DownloadImage(img, filePath)
					if err != nil {
						error = err
					}
					time.Sleep(1 * time.Second)
				}
				chapter.progress = float32(i+1) / float32(len(imgList)) * 100
				progress := 0.0
				for _, c := range chapters {
					progress += float64(c.progress)
				}
				d.Progress = float32(progress / float64(len(chapters)))
				processSend()
			}(i, imgUrl)
		}
	}

	wg.Wait()

	if error != nil {
		d.Fail = true
		return error
	}

	if d.config.PackageType == "cbz" {
		comicInfo := ComicInfo{
			Series:    d.BookInfo.Title,
			Writer:    strings.Join(d.BookInfo.Author, ", "),
			Summary:   d.BookInfo.Description,
			Genre:     strings.Join(d.BookInfo.Genre, ", "),
			Title:     d.Volume.Title,
			Volume:    fmt.Sprintf("%d", d.Index+1),
			PageCount: fmt.Sprintf("%d", total),
		}
		comicInfo.Build(folderPath)
		zipPath := folderPath + ".cbz"
		err := CreateZipFromDirectory(folderPath, zipPath)
		if err != nil {
			return err
		}
		os.RemoveAll(folderPath)
	} else if d.config.PackageType == "zip" {
		zipPath := folderPath + ".zip"
		err := CreateZipFromDirectory(folderPath, zipPath)
		if err != nil {
			return err
		}
		os.RemoveAll(folderPath)
	} else if d.config.PackageType == "epub" {
		zipPath := folderPath + ".epub"
		index := d.Index + 1
		creator := strings.Join(d.BookInfo.Author, ", ")
		MetaData := MetaData{
			Title:       d.Volume.Title,
			Creator:     &creator,
			Description: &d.BookInfo.Description,
			Subject:     d.BookInfo.Genre,
			Index:       &index,
			Series:      &d.BookInfo.Title,
		}

		epubBuilder := EpubBuilder{
			metadata: MetaData,
		}
		epubBuilder.BuildComic(zipPath, folderPath)
		os.RemoveAll(folderPath)
	}
	println(folderPath)

	return nil
}

func (d *DownloaderSingle) DownloadImage(url, filePath string) error {
	maxRetries := 50
	ext := filepath.Ext(filePath)

	for i := 0; i < maxRetries; i++ {
		img, err := GetImage(url)
		if err != nil {
			fmt.Println("Error downloading image:", err)
			d.messageSend(fmt.Sprintf("Error downloading image: %v", err))
			time.Sleep(3 * time.Second)
			continue
		}

		if !isImage(img) {
			fmt.Println("Image is not valid, retrying...")
			d.messageSend("Image is not valid, retrying...")
			time.Sleep(3 * time.Second)
			continue
		}

		if d.config.ImageFormat == "png" {
			filePath = strings.TrimSuffix(filePath, ext) + ".png"
			img, err = ImgToPng(img)
			if err != nil {
				fmt.Println("Error converting image to png:", err)
				d.messageSend(fmt.Sprintf("Error converting image to png: %v", err))
				time.Sleep(3 * time.Second)
				continue
			}
		} else if d.config.ImageFormat == "jpg" {
			filePath = strings.TrimSuffix(filePath, ext) + ".jpg"
			img, err = ImgToJpg(img)
			if err != nil {
				fmt.Println("Error converting image to jpg:", err)
				d.messageSend(fmt.Sprintf("Error converting image to jpg: %v", err))
				time.Sleep(3 * time.Second)
				continue
			}
		}

		os.WriteFile(filePath, img, os.ModePerm)
		if err != nil {
			fmt.Println("Error writing image file:", err)
			d.messageSend(fmt.Sprintf("Error writing image file: %v", err))
			time.Sleep(3 * time.Second)
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to download image: %s", url)
}

func GetNextUrl(html string) (string, error) {
	re := regexp.MustCompile(`url_next:'(.+?)'`)

	match := re.FindStringSubmatch(html)

	if len(match) > 1 {
		return match[1], nil
	} else {
		return "", fmt.Errorf("failed to find next url")
	}
}

func GetPreViousUrl(html string) (string, error) {
	re := regexp.MustCompile(`url_previous:'(.+?)'`)

	match := re.FindStringSubmatch(html)

	if len(match) > 1 {
		return match[1], nil
	} else {
		return "", fmt.Errorf("failed to find previous url")
	}
}

func GetStartUrl(volume *Volume) (string, error) {
	if strings.Contains(volume.Chapters[0].Url, "javascript") {
		if len(volume.Chapters) > 1 {
			html, err := GetText(volume.Chapters[1].Url)
			if err != nil {
				return "", err
			}
			return GetPreViousUrl(html)
		}
		return "", fmt.Errorf("no chapters found")
	}
	return volume.Chapters[0].Url, nil
}
