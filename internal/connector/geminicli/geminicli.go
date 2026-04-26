// Package geminicli discovers MCP servers and skills from a Gemini CLI
// installation by reading its configuration files directly.
//
// Gemini CLI stores config under ~/.gemini/:
//   - ~/.gemini/settings.json              — general settings
//   - ~/.gemini/antigravity/mcp_config.json — MCP server declarations
//   - ~/.gemini/antigravity/skills/         — skill files
//   - <cwd>/.mcp.json                       — project-level MCP config
package geminicli

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"OwnDeck/internal/discovery"
	"OwnDeck/internal/platform"
	"OwnDeck/internal/repository/config"
)

const (
	id   = "gemini-cli"
	name = "Gemini CLI"
)

type Connector struct {
	cfgPaths  []string
	cfgSkills []string
}

func New() *Connector { return &Connector{} }

func NewWithConfig(ac config.AgentConfig) *Connector {
	return &Connector{
		cfgPaths:  ac.ConfigPaths,
		cfgSkills: ac.SkillRoots,
	}
}

func (Connector) ID() string   { return id }
func (Connector) Name() string { return name }

func (c *Connector) Probe() discovery.ClientInfo {
	executable := platform.LookPath("gemini")
	configPaths := platform.ExistingPaths(c.configCandidates()...)
	detected := executable != "" || len(configPaths) > 0

	return discovery.ClientInfo{
		ID:          id,
		Name:        name,
		Detected:    detected,
		Executable:  executable,
		ConfigPaths: configPaths,
		Status:      discovery.ClientStatus(detected),
	}
}

func (c *Connector) DiscoverMCP(_ context.Context) ([]discovery.MCPServer, error) {
	paths := platform.ExistingPaths(c.configCandidates()...)
	if len(paths) == 0 {
		return nil, errors.New("no Gemini CLI MCP config files found")
	}

	var servers []discovery.MCPServer
	for _, p := range paths {
		items, err := platform.ParseMCPServersFile(p)
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
			})
		}
	}
	return servers, nil
}

func (c *Connector) DiscoverSkills(_ context.Context) ([]discovery.SkillAsset, error) {
	skills := platform.DiscoverSkillFiles(c.skillRoots())
	out := make([]discovery.SkillAsset, 0, len(skills))
	for _, s := range skills {
		out = append(out, discovery.SkillAsset{
			Name:        s.Name,
			ClientID:    id,
			Client:      name,
			Description: s.Description,
			SourcePath:  s.Path,
			Scope:       s.Scope,
		})
	}
	return out, nil
}

func (c *Connector) configCandidates() []string {
	if len(c.cfgPaths) > 0 {
		return c.cfgPaths
	}
	home := platform.HomeDir()
	cwd, _ := os.Getwd()
	return []string{
		filepath.Join(home, ".gemini", "settings.json"),
		filepath.Join(home, ".gemini", "antigravity", "mcp_config.json"),
		filepath.Join(cwd, ".mcp.json"),
		filepath.Join(cwd, ".vscode", "mcp.json"),
	}
}

func (c *Connector) skillRoots() []string {
	if len(c.cfgSkills) > 0 {
		return c.cfgSkills
	}
	return []string{
		filepath.Join(platform.HomeDir(), ".gemini", "antigravity", "skills"),
	}
}
