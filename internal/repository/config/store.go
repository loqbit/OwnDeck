// Package config holds OwnDeck's persisted local configuration:
// the AppConfig wire type and a Store interface for swapping
// implementations (file today, SQLite later).
package config

import "errors"

var ErrStoreUnavailable = errors.New("config store unavailable")

type ClientConnection struct {
	Connected   bool   `json:"connected"`
	Permission  string `json:"permission"`
	ConnectedAt string `json:"connectedAt,omitempty"`
}

// AgentConfig stores discovered paths for an AI agent.
// Populated by the first-use scan and persisted so future
// startups skip the filesystem crawl.
type AgentConfig struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Executable  string   `json:"executable,omitempty"`
	ConfigPaths []string `json:"configPaths"`
	SkillRoots  []string `json:"skillRoots"`
	Detected    bool     `json:"detected"`
	ScannedAt   string   `json:"scannedAt"`
}

type AppConfig struct {
	Version int                         `json:"version"`
	Clients map[string]ClientConnection `json:"clients"`
	Agents  []AgentConfig               `json:"agents,omitempty"`
}

type Store interface {
	Path() string
	Load() (AppConfig, error)
	Save(cfg AppConfig) error
}

func ConnectedClientIDs(cfg AppConfig) []string {
	ids := make([]string, 0, len(cfg.Clients))
	for id, conn := range cfg.Clients {
		if conn.Connected {
			ids = append(ids, id)
		}
	}
	return ids
}

// FindAgent returns the AgentConfig for the given agent ID, if present.
func FindAgent(cfg AppConfig, agentID string) (AgentConfig, bool) {
	for _, a := range cfg.Agents {
		if a.ID == agentID {
			return a, true
		}
	}
	return AgentConfig{}, false
}

func defaultConfig() AppConfig {
	return AppConfig{
		Version: 1,
		Clients: map[string]ClientConnection{},
	}
}
