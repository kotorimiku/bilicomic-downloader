package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	bd "bilicomic-downloader/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	donwloaderManager := &bd.DownloaderManager{}
	config := bd.ConfigInstance

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "bilicomic-downloader",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: donwloaderManager.Startup,
		Bind: []interface{}{
			donwloaderManager,
			config,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
