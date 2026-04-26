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
)

const (
	id   = "claude-code"
	name = "Claude Code"
)

type Connector struct{}

func New() *Connector { return &Connector{} }

func (Connector) ID() string   { return id }
func (Connector) Name() string { return name }

func (Connector) Probe() discovery.ClientInfo {
	configPaths := platform.ExistingPaths(configCandidates()...)
	detected := len(configPaths) > 0 || platform.PathExists("/Applications/Claude.app")

	return discovery.ClientInfo{
		ID:          id,
		Name:        name,
		Detected:    detected,
		ConfigPaths: configPaths,
		Status:      discovery.ClientStatus(detected),
	}
}

func (Connector) DiscoverMCP(_ context.Context) ([]discovery.MCPServer, error) {
	paths := platform.ExistingPaths(configCandidates()...)
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

func (Connector) DiscoverSkills(_ context.Context) ([]discovery.SkillAsset, error) {
	return nil, nil
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

func configCandidates() []string {
	home := platform.HomeDir()
	cwd, _ := os.Getwd()
	return []string{
		filepath.Join(home, ".claude", "settings.json"),
		filepath.Join(home, ".claude.json"),
		filepath.Join(cwd, ".mcp.json"),
	}
}
