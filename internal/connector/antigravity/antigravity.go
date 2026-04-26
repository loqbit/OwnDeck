package antigravity

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
	id   = "antigravity"
	name = "Antigravity"
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
	executable := platform.LookPath("antigravity")
	configPaths := platform.ExistingPaths(c.configCandidates()...)
	detected := executable != "" || len(configPaths) > 0 || platform.PathExists("/Applications/Antigravity.app")

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
		return nil, errors.New("no Antigravity MCP config files found")
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
		filepath.Join(home, ".gemini", "antigravity", "mcp_config.json"),
		filepath.Join(home, "Library", "Application Support", "Antigravity", "User", "settings.json"),
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
