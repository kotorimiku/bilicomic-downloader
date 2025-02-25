package bilicomicdownloader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DownloaderSingle struct {
	BookInfo *BookInfo
	Volume   *Volume
	Index    int
	config   *Config
	Progress float32
}

func NewDownloaderSingle(bookInfo *BookInfo, volume *Volume, index int, config *Config) *DownloaderSingle {
	return &DownloaderSingle{
		BookInfo: bookInfo,
		Volume:   volume,
		Index:    index,
		config:   config,
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

func (d *DownloaderSingle) Download(processSend func()) error {
	chapters := d.Volume.Chapters
	folderPath := filepath.Join(d.config.OutputPath, sanitizeFilename(d.BookInfo.Title), sanitizeFilename(d.Volume.Title))
	os.MkdirAll(folderPath, os.ModePerm)
	total := 0

	wg := sync.WaitGroup{}
	maxConcurrency := 1
	sem := make(chan struct{}, maxConcurrency)
	var error error

	for _, chapter := range chapters {
		url := fmt.Sprintf("https://%s%s", d.config.UrlBase, chapter.Url)
		imgList, err := d.GetImageUrlList(url)
		if err != nil {
			return err
		}

		total += len(imgList)
		chapterPath := filepath.Join(folderPath, sanitizeFilename(chapter.Title))
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
				err := DownloadImage(img, filePath)
				if err != nil {
					error = err
				}
				chapter.progress = float32(i+1) / float32(len(imgList)) * 100
				progress := 0.0
				for _, c := range chapters {
					progress += float64(c.progress)
				}
				d.Progress = float32(progress / float64(len(chapters)))
				processSend()
				time.Sleep(1 * time.Second)
			}(i, imgUrl)
		}
	}

	wg.Wait()

	if error != nil {
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

func DownloadImage(url, filePath string) error {
	maxRetries := 50

	for i := 0; i < maxRetries; i++ {
		img, err := GetImage(url)

		if err != nil {
			fmt.Println("Error downloading image:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if !isImage(img) {
			fmt.Println(string(img))
			fmt.Println("Image is not valid")
			time.Sleep(3 * time.Second)
			continue
		}

		err = os.WriteFile(filePath, img, os.ModePerm)
		if err != nil {
			fmt.Println("Error writing image file:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to download image: %s", url)
}
