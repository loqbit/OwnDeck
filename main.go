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
	"OwnDeck/internal/connector/generic"
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

	// Build connectors from ALL discovered agents.
	// Known agents get specialized connectors; unknown ones get the generic adapter.
	registry := connector.NewRegistry(buildConnectors(cfg)...)

	discoverySvc := discoverysvc.New(registry)
	connectionSvc := connectionsvc.New(store)

	return NewApp(discoverySvc, connectionSvc, store), nil
}

// specializedConnectors maps agent IDs to their dedicated connector
// constructors. Agents not in this map use the generic adapter.
var specializedConnectors = map[string]func(config.AgentConfig) connector.Connector{
	"codex":          func(ac config.AgentConfig) connector.Connector { return codex.NewWithConfig(ac) },
	"claude-code":    func(ac config.AgentConfig) connector.Connector { return claudecode.NewWithConfig(ac) },
	"claude-desktop": func(ac config.AgentConfig) connector.Connector { return claudedesktop.NewWithConfig(ac) },
	"antigravity":    func(ac config.AgentConfig) connector.Connector { return antigravity.NewWithConfig(ac) },
}

// buildConnectors creates a connector for every discovered agent.
// Known agents get their specialized connector; everything else
// gets the generic JSON-based adapter.
func buildConnectors(cfg config.AppConfig) []connector.Connector {
	var connectors []connector.Connector

	for _, agent := range cfg.Agents {
		if !agent.Detected {
			continue // skip agents that weren't actually found
		}
		if factory, ok := specializedConnectors[agent.ID]; ok {
			connectors = append(connectors, factory(agent))
		} else {
			connectors = append(connectors, generic.NewFromAgentConfig(agent))
		}
	}

	return connectors
}
