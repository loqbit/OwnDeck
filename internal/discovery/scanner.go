// Package discovery defines the wire types returned by the Wails layer
// and the system scanner that probes for installed AI agents.
package discovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"OwnDeck/internal/platform"
	"OwnDeck/internal/repository/config"
)

// AgentProbe describes how to detect a single AI agent on the system.
type AgentProbe struct {
	ID              string   // e.g. "claude-code"
	Name            string   // e.g. "Claude Code"
	ExecutableNames []string // binaries to search in PATH
	AppBundlePaths  []string // macOS .app bundles to check
	ConfigPatterns  []string // paths to check for MCP config files
	SkillPatterns   []string // directories to scan for skills
}

// builtinProbes lists the known AI agents and where to look for them
// on macOS. This is the "seed" data — results are persisted so we
// don't re-scan on every launch.
func builtinProbes() []AgentProbe {
	home := platform.HomeDir()
	cwd, _ := os.Getwd()

	return []AgentProbe{
		{
			ID:              "claude-code",
			Name:            "Claude Code",
			ExecutableNames: []string{"claude"},
			AppBundlePaths:  []string{"/Applications/Claude.app"},
			ConfigPatterns: []string{
				filepath.Join(home, ".claude", "settings.json"),
				filepath.Join(home, ".claude.json"),
				filepath.Join(cwd, ".mcp.json"),
			},
			SkillPatterns: []string{
				filepath.Join(home, ".claude", "skills"),
			},
		},
		{
			ID:              "claude-desktop",
			Name:            "Claude Desktop",
			ExecutableNames: nil,
			AppBundlePaths:  []string{"/Applications/Claude.app"},
			ConfigPatterns: []string{
				filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json"),
			},
			SkillPatterns: nil,
		},
		{
			ID:              "codex",
			Name:            "Codex",
			ExecutableNames: []string{"codex"},
			AppBundlePaths:  nil,
			ConfigPatterns: []string{
				filepath.Join(home, ".codex", "config.toml"),
			},
			SkillPatterns: []string{
				filepath.Join(home, ".codex", "skills"),
				filepath.Join(home, ".codex", "plugins", "cache"),
			},
		},
		{
			ID:              "gemini-cli",
			Name:            "Gemini CLI",
			ExecutableNames: []string{"gemini"},
			AppBundlePaths:  nil,
			ConfigPatterns: []string{
				filepath.Join(home, ".gemini", "settings.json"),
				filepath.Join(home, ".gemini", "antigravity", "mcp_config.json"),
				filepath.Join(cwd, ".mcp.json"),
				filepath.Join(cwd, ".vscode", "mcp.json"),
			},
			SkillPatterns: []string{
				filepath.Join(home, ".gemini", "antigravity", "skills"),
			},
		},
		{
			ID:              "antigravity",
			Name:            "Antigravity",
			ExecutableNames: []string{"antigravity"},
			AppBundlePaths:  []string{"/Applications/Antigravity.app"},
			ConfigPatterns:  antigravityConfigPatterns(home),
			SkillPatterns:   nil,
		},
	}
}

func antigravityConfigPatterns(home string) []string {
	paths := []string{
		filepath.Join(home, "Library", "Application Support", "Antigravity", "User", "settings.json"),
	}

	// Antigravity extensions may bundle their own .mcp.json files.
	// By explicitly listing them in ConfigPatterns, we ensure they belong to
	// the "Antigravity" agent rather than showing up as "Unknown" agents.
	extDir := filepath.Join(home, ".antigravity", "extensions")
	if entries, err := os.ReadDir(extDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				mcpPath := filepath.Join(extDir, entry.Name(), ".mcp.json")
				if platform.PathExists(mcpPath) {
					paths = append(paths, mcpPath)
				}
			}
		}
	}
	return paths
}

// ScanAgents probes the local filesystem for known AI agents and
// returns a configuration entry for each one found.
// This is designed to run once on first use; results are persisted
// in AppConfig.Agents so subsequent launches skip the scan.
func ScanAgents() []config.AgentConfig {
	now := time.Now().UTC().Format(time.RFC3339)
	var agents []config.AgentConfig

	for _, probe := range builtinProbes() {
		agent := probeAgent(probe, now)
		agents = append(agents, agent)
	}

	// Also try to discover unknown agents by scanning common
	// config directories for files containing "mcpServers".
	agents = append(agents, scanUnknownAgents(now)...)

	return agents
}

func probeAgent(probe AgentProbe, now string) config.AgentConfig {
	agent := config.AgentConfig{
		ID:        probe.ID,
		Name:      probe.Name,
		ScannedAt: now,
	}

	// Check executable
	for _, exe := range probe.ExecutableNames {
		if path := platform.LookPath(exe); path != "" {
			agent.Executable = path
			agent.Detected = true
			break
		}
	}

	// Check .app bundles (macOS)
	if !agent.Detected {
		for _, app := range probe.AppBundlePaths {
			if platform.PathExists(app) {
				agent.Detected = true
				break
			}
		}
	}

	// Collect existing config paths
	agent.ConfigPaths = platform.ExistingPaths(probe.ConfigPatterns...)
	if len(agent.ConfigPaths) > 0 {
		agent.Detected = true
	}

	// Collect existing skill roots
	agent.SkillRoots = platform.ExistingPaths(probe.SkillPatterns...)

	return agent
}

// scanUnknownAgents crawls common config directories looking for
// JSON files that contain an "mcpServers" key. Any file found that
// doesn't belong to a known probe is registered as an unknown agent.
func scanUnknownAgents(now string) []config.AgentConfig {
	home := platform.HomeDir()
	if home == "" {
		return nil
	}

	knownPaths := map[string]bool{}
	for _, p := range builtinProbes() {
		for _, cp := range p.ConfigPatterns {
			knownPaths[cp] = true
		}
	}

	searchDirs := []string{
		home, // scans ~/.* dot-config dirs
		filepath.Join(home, "Library", "Application Support"),
	}

	var agents []config.AgentConfig
	seen := map[string]bool{}

	for _, dir := range searchDirs {
		_ = filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return filepath.SkipDir
			}

			// Depth limit: only go 3 levels deep from the search root
			rel, _ := filepath.Rel(dir, path)
			depth := strings.Count(rel, string(filepath.Separator))
			if entry.IsDir() && depth >= 3 {
				return filepath.SkipDir
			}

			// Skip non-JSON files
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				return nil
			}

			// Skip known paths and already-seen files
			if knownPaths[path] || seen[path] {
				return nil
			}

			// Check if it contains "mcpServers"
			if hasMCPServersKey(path) {
				seen[path] = true
				parentDir := filepath.Base(filepath.Dir(path))
				agents = append(agents, config.AgentConfig{
					ID:          "unknown-" + parentDir,
					Name:        "Unknown (" + parentDir + ")",
					ConfigPaths: []string{path},
					Detected:    true,
					ScannedAt:   now,
				})
			}
			return nil
		})
	}

	return agents
}

// hasMCPServersKey does a fast check whether a JSON file contains
// an "mcpServers" top-level key.
func hasMCPServersKey(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil || len(data) > 10*1024*1024 { // skip files >10MB
		return false
	}

	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return false
	}
	_, ok := root["mcpServers"]
	return ok
}
