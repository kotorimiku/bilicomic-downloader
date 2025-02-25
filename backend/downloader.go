package bilicomicdownloader

import (
	"fmt"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Downloader struct {
	bookId   string
	bookInfo *BookInfo
	Volumes  []*Volume
	config   *Config
}

func NewDownloader(bookId string, config *Config) *Downloader {
	return &Downloader{
		bookId:   bookId,
		config:   config,
		bookInfo: &BookInfo{},
	}
}

func (d *Downloader) GetDownloadList(chapters []int) []*DownloaderSingle {
	var downloaderSinglesList []*DownloaderSingle = make([]*DownloaderSingle, 0, len(chapters))
	for _, index := range chapters {
		downloaderSingle := DownloaderSingle{
			Volume:   d.Volumes[index],
			BookInfo: d.bookInfo,
			config:   d.config,
			Index:    index,
		}
		downloaderSinglesList = append(downloaderSinglesList, &downloaderSingle)
	}
	return downloaderSinglesList
}

func (d *Downloader) DownloadList(volumes []int) error {
	for _, index := range volumes {
		volume := d.Volumes[index]
		downloadSingle := NewDownloaderSingle(d.bookInfo, volume, index, d.config)
		err := downloadSingle.Download(func() {})
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Downloader) GetMetadata() error {
	text, err := GetRaw(fmt.Sprintf("https://%s/detail/%s.html", d.config.UrlBase, d.bookId))
	if err != nil {
		return err
	}

	html, err := goquery.NewDocumentFromReader(text)

	if err != nil {
		return err
	}

	book_info := html.Find("div.book-detail-info").First()

	title := book_info.Find("h1.book-title").First().Text()

	author := make([]string, 0, 2)
	book_author := book_info.Find("div.book-rand-a").First()
	book_author.Find("a").Each(func(i int, s *goquery.Selection) {
		author = append(author, s.Text())
	})

	genre := make([]string, 0, 10)
	book_info.Find("p.book-meta span em").Each(func(i int, s *goquery.Selection) {
		genre = append(genre, s.Text())
	})

	bookSummary := html.Find("#bookSummary").First()
	description := bookSummary.Find("content").First().Text()

	cover := book_info.Find("img").First().AttrOr("src", "")

	d.bookInfo = &BookInfo{
		Title:       title,
		Author:      author,
		Genre:       genre,
		Description: description,
		Cover:       cover,
	}

	return nil
}

func (d *Downloader) GetVolume() error {
	text, err := GetRaw(fmt.Sprintf("https://%s/read/%s/catalog", d.config.UrlBase, d.bookId))
	if err != nil {
		return err
	}

	html, err := goquery.NewDocumentFromReader(text)
	if err != nil {
		return err
	}

	volumes := []*Volume{}

	volumesNode := html.Find("#volumes").First()
	volumesNode.Find("div.catalog-volume").Each(func(i int, s *goquery.Selection) {
		vol_title := s.Find("li.chapter-bar.chapter-li").First().Text()
		cover := s.Find("img").First().AttrOr("src", "")
		var chapters = make([]*Chapter, 0, 15)
		s.Find("li.chapter-li.jsChapter").Each(func(i int, s *goquery.Selection) {
			chapter_title := s.Text()
			chapter_url := s.Find("a").First().AttrOr("href", "")
			chapters = append(chapters, &Chapter{
				Title: chapter_title,
				Url:   chapter_url,
			})
		})
		volumes = append(volumes, &Volume{
			Title:    vol_title,
			Cover:    cover,
			Chapters: chapters,
		})
	})
	d.Volumes = volumes
	return nil
}

var maxConcurrency = 1
var sem = make(chan struct{}, maxConcurrency)
var isDownloading = false

func DownloadList(downloaderSingleList chan *DownloaderSingle, processSend func(), clearDownloaders func()) {
	if isDownloading {
		return
	}
	isDownloading = true
	var wg sync.WaitGroup

	for downloaderSingle := range downloaderSingleList {
		wg.Add(1)
		sem <- struct{}{}
		go func(downloaderSingle *DownloaderSingle) {
			defer wg.Done()
			err := downloaderSingle.Download(processSend)
			if err != nil {
				fmt.Println("Error downloading chapter:", err)
			}
			clearDownloaders()
			<-sem
		}(downloaderSingle)
	}
	wg.Wait()
	isDownloading = false
}
