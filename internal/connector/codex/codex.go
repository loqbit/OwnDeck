// Package codex discovers MCP servers from a Codex installation by
// reading ~/.codex/config.toml directly. It aggregates two distinct
// sources that the Codex desktop UI is currently known to merge
// incorrectly (see openai/codex#17360):
//
//  1. User-defined servers under [mcp_servers.<name>] in config.toml.
//  2. Plugin-bundled servers shipped inside enabled plugin directories
//     under ~/.codex/plugins/cache/<marketplace>/<plugin>/<version>/,
//     each pointed to by the plugin manifest's "mcpServers" field.
//
// Each returned MCPServer is tagged with its Origin so the UI can
// distinguish "user" entries from "plugin:<id>" ones.
package codex

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"

	"OwnDeck/internal/discovery"
	"OwnDeck/internal/platform"
)

const (
	id   = "codex"
	name = "Codex"
)

type Connector struct{}

func New() *Connector { return &Connector{} }

func (Connector) ID() string   { return id }
func (Connector) Name() string { return name }

func (c Connector) Probe() discovery.ClientInfo {
	configPaths := platform.ExistingPaths(configPath())
	detected := len(configPaths) > 0 || platform.PathExists(codexHome())
	return discovery.ClientInfo{
		ID:          id,
		Name:        name,
		Detected:    detected,
		ConfigPaths: configPaths,
		Status:      discovery.ClientStatus(detected),
	}
}

// DiscoverMCP reads config.toml plus every enabled plugin's bundled
// .mcp.json. Errors from individual plugins are swallowed so one
// broken plugin doesn't block discovery of the rest.
func (Connector) DiscoverMCP(_ context.Context) ([]discovery.MCPServer, error) {
	cfg, err := loadConfig(configPath())
	if err != nil {
		return nil, err
	}

	servers := userServers(cfg, configPath())
	servers = append(servers, pluginServers(cfg)...)
	return servers, nil
}

func (Connector) DiscoverSkills(_ context.Context) ([]discovery.SkillAsset, error) {
	skills := platform.DiscoverSkillFiles(skillRoots())
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

// codexConfig models only the parts of ~/.codex/config.toml we care
// about. Unknown fields in the file are ignored by toml.
type codexConfig struct {
	MCPServers map[string]tomlMCPServer       `toml:"mcp_servers"`
	Plugins    map[string]tomlPluginEntry     `toml:"plugins"`
	Markets    map[string]tomlMarketplaceMeta `toml:"marketplaces"`
}

type tomlMCPServer struct {
	Command string            `toml:"command"`
	Args    []string          `toml:"args"`
	URL     string            `toml:"url"`
	Cwd     string            `toml:"cwd"`
	Env     map[string]string `toml:"env"`
	Auth    string            `toml:"auth"`
}

type tomlPluginEntry struct {
	Enabled bool `toml:"enabled"`
}

type tomlMarketplaceMeta struct {
	Source     string `toml:"source"`
	SourceType string `toml:"source_type"`
}

func loadConfig(path string) (codexConfig, error) {
	if !platform.PathExists(path) {
		return codexConfig{}, errors.New("codex config.toml not found at " + path)
	}
	var cfg codexConfig
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return codexConfig{}, err
	}
	return cfg, nil
}

func userServers(cfg codexConfig, sourcePath string) []discovery.MCPServer {
	names := sortedKeys(cfg.MCPServers)
	out := make([]discovery.MCPServer, 0, len(names))
	for _, n := range names {
		raw := cfg.MCPServers[n]
		out = append(out, discovery.MCPServer{
			Name:       n,
			ClientID:   id,
			Client:     name,
			Transport:  inferTransport(raw.URL, raw.Command),
			Command:    raw.Command,
			Args:       strings.Join(raw.Args, " "),
			URL:        raw.URL,
			Env:        envSummary(raw.Env),
			Cwd:        raw.Cwd,
			Auth:       raw.Auth,
			Status:     "configured",
			SourcePath: sourcePath,
			Origin:     "user",
			OriginPath: sourcePath,
		})
	}
	return out
}

// pluginServers walks every enabled plugin directory and pulls MCP
// servers out of the plugin's manifest-referenced .mcp.json file.
// This is the bit Codex's own desktop UI currently misses (#17360).
func pluginServers(cfg codexConfig) []discovery.MCPServer {
	pluginIDs := enabledPluginIDs(cfg)
	if len(pluginIDs) == 0 {
		return nil
	}

	var out []discovery.MCPServer
	for _, pluginID := range pluginIDs {
		dir := pluginDir(pluginID)
		if dir == "" {
			continue
		}
		out = append(out, readPluginMCP(pluginID, dir)...)
	}
	return out
}

// enabledPluginIDs returns plugin IDs (e.g. "github@openai-curated")
// for which config.toml has [plugins."NAME@MARKETPLACE"] enabled = true.
func enabledPluginIDs(cfg codexConfig) []string {
	ids := make([]string, 0, len(cfg.Plugins))
	for pluginID, entry := range cfg.Plugins {
		if entry.Enabled {
			ids = append(ids, pluginID)
		}
	}
	sort.Strings(ids)
	return ids
}

// pluginDir resolves a plugin ID like "github@openai-curated" to the
// most recent versioned directory under ~/.codex/plugins/cache/.
// Plugins live at <cache>/<marketplace>/<plugin>/<version>/.
func pluginDir(pluginID string) string {
	plugin, market, ok := strings.Cut(pluginID, "@")
	if !ok || plugin == "" || market == "" {
		return ""
	}
	base := filepath.Join(codexHome(), "plugins", "cache", market, plugin)
	entries, err := os.ReadDir(base)
	if err != nil {
		return ""
	}
	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}
	if len(versions) == 0 {
		return ""
	}
	// We don't try to parse semver; lexical sort is good enough to
	// pick a stable directory when multiple versions are cached.
	sort.Strings(versions)
	return filepath.Join(base, versions[len(versions)-1])
}

// readPluginMCP loads the plugin manifest, follows its mcpServers
// pointer (a relative path), and parses the resulting .mcp.json.
func readPluginMCP(pluginID string, pluginDir string) []discovery.MCPServer {
	manifestPath := filepath.Join(pluginDir, ".codex-plugin", "plugin.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil
	}

	var manifest struct {
		MCPServers any `json:"mcpServers"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil
	}

	mcpPath, ok := manifest.MCPServers.(string)
	if !ok || mcpPath == "" {
		// Some plugins inline mcpServers as an object instead of a
		// path. We don't handle that yet — skip and leave it for a
		// future iteration.
		return nil
	}

	resolved := filepath.Join(pluginDir, mcpPath)
	items, err := platform.ParseMCPServersFile(resolved)
	if err != nil {
		return nil
	}

	out := make([]discovery.MCPServer, 0, len(items))
	for _, item := range items {
		out = append(out, discovery.MCPServer{
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
			Origin:     "plugin:" + pluginID,
			OriginPath: pluginDir,
		})
	}
	return out
}

func configPath() string {
	return filepath.Join(codexHome(), "config.toml")
}

func codexHome() string {
	home := platform.HomeDir()
	if home == "" {
		return ""
	}
	return filepath.Join(home, ".codex")
}

func skillRoots() []string {
	return []string{
		filepath.Join(codexHome(), "skills"),
		filepath.Join(codexHome(), "plugins", "cache"),
	}
}

func sortedKeys(m map[string]tomlMCPServer) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func envSummary(env map[string]string) string {
	if len(env) == 0 {
		return "-"
	}
	return "set"
}

func inferTransport(url, command string) string {
	if url != "" {
		return "http"
	}
	if command != "" {
		return "stdio"
	}
	return "unknown"
}
