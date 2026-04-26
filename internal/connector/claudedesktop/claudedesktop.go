// Package claudedesktop discovers MCP servers from the Claude
// Desktop macOS app by reading claude_desktop_config.json directly.
//
// Config location (macOS):
//
//	~/Library/Application Support/Claude/claude_desktop_config.json
//
// The file may or may not contain an "mcpServers" key; users who
// have not added any MCP servers will simply have a "preferences"
// blob and nothing else. We treat that as "detected, zero servers".
package claudedesktop

import (
	"context"
	"errors"
	"path/filepath"

	"OwnDeck/internal/discovery"
	"OwnDeck/internal/platform"
	"OwnDeck/internal/repository/config"
)

const (
	id   = "claude-desktop"
	name = "Claude Desktop"
)

type Connector struct {
	cfgPaths []string
}

func New() *Connector { return &Connector{} }

func NewWithConfig(ac config.AgentConfig) *Connector {
	return &Connector{cfgPaths: ac.ConfigPaths}
}

func (Connector) ID() string   { return id }
func (Connector) Name() string { return name }

func (c *Connector) Probe() discovery.ClientInfo {
	configPaths := platform.ExistingPaths(c.configCandidates()...)
	detected := len(configPaths) > 0 || platform.PathExists("/Applications/Claude.app")

	return discovery.ClientInfo{
		ID:          id,
		Name:        name,
		Detected:    detected,
		ConfigPaths: configPaths,
		Status:      discovery.ClientStatus(detected),
	}
}

func (c *Connector) DiscoverMCP(_ context.Context) ([]discovery.MCPServer, error) {
	paths := platform.ExistingPaths(c.configCandidates()...)
	if len(paths) == 0 {
		return nil, errors.New("no Claude Desktop config files found")
	}

	var servers []discovery.MCPServer
	for _, path := range paths {
		items, err := platform.ParseMCPServersFile(path)
		if err != nil {
			continue
		}
		for _, item := range items {
			servers = append(servers, discovery.MCPServer{
				Name:       item.Name,
				ClientID:   id,
				Client:     name,
				Transport:  item.Transport,
				Command:    item.Command,
				Args:       item.Args,
				URL:        item.URL,
				Env:        item.Env,
				Cwd:        item.Cwd,
				Status:     "configured",
				SourcePath: item.Source,
				Origin:     "user",
				OriginPath: path,
			})
		}
	}
	return servers, nil
}

// Claude Desktop has no Skills concept (it's a chat app, not a
// coding agent), so this is intentionally empty.
func (Connector) DiscoverSkills(_ context.Context) ([]discovery.SkillAsset, error) {
	return nil, nil
}

func (c *Connector) configCandidates() []string {
	if len(c.cfgPaths) > 0 {
		return c.cfgPaths
	}
	home := platform.HomeDir()
	return []string{
		filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json"),
	}
}
