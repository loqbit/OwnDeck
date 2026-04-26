// Package generic provides a connector for any AI agent discovered by
// the system scanner that doesn't have a dedicated connector implementation.
// It reads MCP servers from JSON config files containing "mcpServers"
// and discovers skills from SKILL.md files in the provided roots.
package generic

import (
	"context"
	"errors"

	"OwnDeck/internal/discovery"
	"OwnDeck/internal/platform"
	"OwnDeck/internal/repository/config"
)

// Connector is a generic MCP client adapter that works for any agent
// whose config files use the standard { "mcpServers": { ... } } format.
type Connector struct {
	id         string
	name       string
	configPaths []string
	skillRoots  []string
	detected    bool
	executable  string
}

// NewFromAgentConfig creates a Connector from a scanner-discovered AgentConfig.
func NewFromAgentConfig(ac config.AgentConfig) *Connector {
	return &Connector{
		id:          ac.ID,
		name:        ac.Name,
		configPaths: ac.ConfigPaths,
		skillRoots:  ac.SkillRoots,
		detected:    ac.Detected,
		executable:  ac.Executable,
	}
}

func (c *Connector) ID() string   { return c.id }
func (c *Connector) Name() string { return c.name }

func (c *Connector) Probe() discovery.ClientInfo {
	configPaths := platform.ExistingPaths(c.configPaths...)
	detected := c.detected || len(configPaths) > 0

	return discovery.ClientInfo{
		ID:          c.id,
		Name:        c.name,
		Detected:    detected,
		Executable:  c.executable,
		ConfigPaths: configPaths,
		Status:      discovery.ClientStatus(detected),
	}
}

func (c *Connector) DiscoverMCP(_ context.Context) ([]discovery.MCPServer, error) {
	paths := platform.ExistingPaths(c.configPaths...)
	if len(paths) == 0 {
		return nil, errors.New("no config files found for " + c.name)
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
				ClientID:   c.id,
				Client:     c.name,
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

func (c *Connector) DiscoverSkills(_ context.Context) ([]discovery.SkillAsset, error) {
	if len(c.skillRoots) == 0 {
		return nil, nil
	}
	skills := platform.DiscoverSkillFiles(c.skillRoots)
	out := make([]discovery.SkillAsset, 0, len(skills))
	for _, s := range skills {
		out = append(out, discovery.SkillAsset{
			Name:        s.Name,
			ClientID:    c.id,
			Client:      c.name,
			Description: s.Description,
			SourcePath:  s.Path,
			Scope:       s.Scope,
		})
	}
	return out, nil
}
