// Package claudecode discovers MCP servers from a Claude Code
// installation by reading its configuration files directly.
// We never invoke the `claude` CLI: parsing JSON is faster, more
// faithful, and works whether or not the CLI is on PATH.
//
// Claude Code keeps MCP server declarations in three places:
//
//   - ~/.claude/settings.json   — user-level
//   - ~/.claude.json            — user-level (legacy / unified)
//   - <cwd>/.mcp.json           — project-level
//
// We aggregate them and tag each server with its origin source.
package claudecode

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"OwnDeck/internal/discovery"
	"OwnDeck/internal/platform"
	"OwnDeck/internal/repository/config"
)

const (
	id   = "claude-code"
	name = "Claude Code"
)

// Connector discovers MCP servers and skills from Claude Code.
// When created with NewWithConfig, it uses persisted paths from
// the scanner; otherwise it falls back to hardcoded defaults.
type Connector struct {
	cfgPaths   []string
	cfgSkills  []string
}

// New creates a Connector with hardcoded default paths.
func New() *Connector { return &Connector{} }

// NewWithConfig creates a Connector using previously discovered
// paths from the agent scanner. Falls back to defaults if the
// AgentConfig has empty paths.
func NewWithConfig(ac config.AgentConfig) *Connector {
	return &Connector{
		cfgPaths:  ac.ConfigPaths,
		cfgSkills: ac.SkillRoots,
	}
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
		return nil, errors.New("no Claude Code MCP config files found")
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
				Origin:     originForPath(path),
				OriginPath: path,
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

// originForPath labels a config path as "user" (under the user's home
// dot-config) or "project" (a .mcp.json sitting in a working dir).
func originForPath(path string) string {
	home := platform.HomeDir()
	if home != "" && strings.HasPrefix(path, home) {
		return "user"
	}
	return "project"
}

func (c *Connector) configCandidates() []string {
	if len(c.cfgPaths) > 0 {
		return c.cfgPaths
	}
	return defaultConfigCandidates()
}

func defaultConfigCandidates() []string {
	home := platform.HomeDir()
	cwd, _ := os.Getwd()
	return []string{
		filepath.Join(home, ".claude", "settings.json"),
		filepath.Join(home, ".claude.json"),
		filepath.Join(cwd, ".mcp.json"),
	}
}

func (c *Connector) skillRoots() []string {
	if len(c.cfgSkills) > 0 {
		return c.cfgSkills
	}
	return defaultSkillRoots()
}

func defaultSkillRoots() []string {
	home := platform.HomeDir()
	return []string{
		filepath.Join(home, ".claude", "skills"),
	}
}
