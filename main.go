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
	"OwnDeck/internal/discovery"
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

// buildApp wires the dependency graph: store -> scanner -> connectors -> services -> handler.
func buildApp() (*App, error) {
	store, err := config.NewFileStore()
	if err != nil {
		return nil, err
	}

	// Run agent scan on first use (or when agents list is empty).
	cfg, err := store.Load()
	if err != nil {
		return nil, err
	}
	if len(cfg.Agents) == 0 {
		log.Println("First use: scanning for installed AI agents...")
		cfg.Agents = discovery.ScanAgents()
		if err := store.Save(cfg); err != nil {
			log.Printf("Warning: failed to persist agent scan: %v", err)
		}
		log.Printf("Found %d agents", len(cfg.Agents))
	}

	// Create connectors using discovered paths (falls back to defaults
	// if no AgentConfig exists for a given connector).
	registry := connector.NewRegistry(
		makeConnector("codex", cfg, func(ac config.AgentConfig) connector.Connector { return codex.NewWithConfig(ac) }, func() connector.Connector { return codex.New() }),
		makeConnector("claude-code", cfg, func(ac config.AgentConfig) connector.Connector { return claudecode.NewWithConfig(ac) }, func() connector.Connector { return claudecode.New() }),
		makeConnector("claude-desktop", cfg, func(ac config.AgentConfig) connector.Connector { return claudedesktop.NewWithConfig(ac) }, func() connector.Connector { return claudedesktop.New() }),
		makeConnector("antigravity", cfg, func(ac config.AgentConfig) connector.Connector { return antigravity.NewWithConfig(ac) }, func() connector.Connector { return antigravity.New() }),
	)

	discoverySvc := discoverysvc.New(registry)
	connectionSvc := connectionsvc.New(store)

	return NewApp(discoverySvc, connectionSvc, store), nil
}

// makeConnector creates a connector using the scanned AgentConfig if
// available, otherwise falls back to the zero-config constructor.
func makeConnector(
	agentID string,
	cfg config.AppConfig,
	withConfig func(config.AgentConfig) connector.Connector,
	fallback func() connector.Connector,
) connector.Connector {
	if ac, ok := config.FindAgent(cfg, agentID); ok {
		return withConfig(ac)
	}
	return fallback()
}
