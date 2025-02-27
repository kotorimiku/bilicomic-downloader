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
		ClientInit("night=0; cf_chl_rc_m=36; jieqiRecentRead=78.4616.0.1.1740640932.0; cf_clearance=VP5cUD99obtGMujKzDNsBsmlZJ3t34XiVqu2qmxL5WQ-1740640933-1.2.1.1-C_aX0gLAa3EFOrv6GFkWkE4RK1C2uvRuI1PGU75TztIPd6F50BEOt36aAFzZZRgxgHI_eMwtZbzfx4is9stVkM4t8YSK0z69QAOB3QkTm.yqyAK1YzwoQIkVaDm996Gwq72Zb.UcHzHYEH.pSwrrSNaMoUN_RpSvoOxmL_I5dkNKeGg_O2jHcBQ83j_jXz3wyOOhNwTVTFu1DoWEblr6mW57JFbrwagijL4bqnrYjtjNBlmV.eKp_7bPd1cLoOT.dC_Tpiqxq.RS2fd3lT5Ah7ht2yu94nI.Ps4LinS10K8gUng6wVVFeYD_7QGZ80r.DWTYBvedM9Q5G_te4fq_5Q")
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

func (d *DownloaderManager) GetBookInfo() (*BookInfo, error) {
	err := d.view.GetMetadata()
	if err != nil {
		return nil, err
	}
	return d.view.bookInfo, nil
}

func (d *DownloaderManager) GetChapter() ([]*Volume, error) {
	err := d.view.GetVolume()
	if err != nil {
		return nil, err
	}
	return d.view.Volumes, nil
}

func (d *DownloaderManager) DownloadList(chapters []int) {
	downloaderSingleList := d.view.GetDownloadList(chapters, d.MessageSend)
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

func (d *DownloaderManager) MessageSend(message string) {
	runtime.EventsEmit(d.ctx, "message", message)
}
