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

type AppConfig struct {
	Version int                         `json:"version"`
	Clients map[string]ClientConnection `json:"clients"`
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

func defaultConfig() AppConfig {
	return AppConfig{
		Version: 1,
		Clients: map[string]ClientConnection{},
	}
}
