package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"OwnDeck/internal/connector"
	"OwnDeck/internal/connector/antigravity"
	"OwnDeck/internal/connector/claudecode"
	"OwnDeck/internal/connector/claudedesktop"
	"OwnDeck/internal/connector/codex"
	"OwnDeck/internal/repository/config"
	"OwnDeck/internal/service/connectionsvc"
	"OwnDeck/internal/service/discoverysvc"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app, err := buildApp()
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	err = wails.Run(&options.App{
		Title:  "OwnDeck",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 9, G: 9, B: 11, A: 1},
		OnStartup:        app.startup,
		Mac: &mac.Options{
			Appearance:           mac.DefaultAppearance,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
		},
		Bind: []interface{}{app},
	})
	if err != nil {
		log.Printf("wails: %v", err)
	}
}

// buildApp wires the dependency graph: store -> services -> handler.
// Adding a new connector is a single line in NewRegistry below.
func buildApp() (*App, error) {
	store, err := config.NewFileStore()
	if err != nil {
		return nil, err
	}

	registry := connector.NewRegistry(
		codex.New(),
		claudecode.New(),
		claudedesktop.New(),
		antigravity.New(),
	)

	discoverySvc := discoverysvc.New(registry)
	connectionSvc := connectionsvc.New(store)

	return NewApp(discoverySvc, connectionSvc), nil
}
