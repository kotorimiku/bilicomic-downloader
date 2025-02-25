package bilicomicdownloader

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var muD sync.Mutex

type DownloaderManager struct {
	ctx             context.Context
	downloaders     []*DownloaderSingle
	downloadersChan chan *DownloaderSingle
	view            *Downloader
}

func (d *DownloaderManager) Startup(ctx context.Context) {
	d.ctx = ctx
	if ConfigInstance.Cookie == "" {
		ClientInit("night=0; jieqiRecentRead=60.3373.0.1.1740412640.0-159.13203.0.1.1740416624.0; cf_clearance=6Kz4rMl_t4rbmXSHDAYzhHKxLYobH5HqMMOm1JeuhbI-1740416624-1.2.1.1-tRzR8okVuw2PcYZLyMgRFGruVXkUWl75MXfn1LpWHSiN.aqauXu7HHHqTTVtxEEP9MqIxdhvSO81imXizmR.dteRuXwsZwpBgsJfzxOCg8Oq3xEGBtc0MUQUjqLeWT4BQK2L1U0aCRDnuy9psaz3NnFzvy2_VqDiFUFOmzdPIEC5PXDugPRmhHD86bIsGc49rzxWIL2Vb7tWqNwHb98.6daM8LHb8B5Wj5e2jsUrdX.Iuo_lj9QFgiwQQm80AOF4vSx8cTujxPGu8BAuomsHUTMNBWiBzTiuBbrKJyjRVAK_mYYLHXoMdmbUWlS_h1HTPGwUa.dX9JMMjkY0pH8MqA")
	} else {
		ClientInit(ConfigInstance.Cookie)
	}
	d.downloadersChan = make(chan *DownloaderSingle, 999)
	DownloadList(d.downloadersChan, d.ProcessSend, d.ClearDownloaders)
}

// func (d *DownloaderManager) Search(keyword string, page int) ([]Comic, error) {
// 	return Search(ConfigInstance.UrlBase, keyword, page)
// }

func (d *DownloaderManager) GetDownloader(bookId string) {
	d.view = NewDownloader(bookId, ConfigInstance)
}

func (d *DownloaderManager) GetBookInfo() (BookInfo, error) {
	d.view.GetMetadata()
	return *d.view.bookInfo, nil
}

func (d *DownloaderManager) GetChapter() ([]*Volume, error) {
	d.view.GetVolume()
	return d.view.Volumes, nil
}

func (d *DownloaderManager) DownloadList(chapters []int) {
	downloaderSingleList := d.view.GetDownloadList(chapters)
	for _, downloaderSingle := range downloaderSingleList {
		d.downloadersChan <- downloaderSingle
	}
	muD.Lock()
	d.downloaders = append(d.downloaders, downloaderSingleList...)
	muD.Unlock()
}

func (d *DownloaderManager) GetDownloaders() []*DownloaderSingle {
	return d.downloaders
}

func (d *DownloaderManager) ClearDownloaders() {
	index := 0
	muD.Lock()
	for i, downloader := range d.downloaders {
		if downloader.Progress < 100 {
			break
		}
		index = i + 1
	}
	d.downloaders = d.downloaders[index:]
	muD.Unlock()
	d.ProcessSend()
}

func (d *DownloaderManager) ProcessSend() {
	runtime.EventsEmit(d.ctx, "progress", d.downloaders)
}
