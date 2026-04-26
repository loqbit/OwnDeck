package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type FileStore struct {
	path string
}

func NewFileStore() (*FileStore, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return &FileStore{
		path: filepath.Join(baseDir, "OwnDeck", "config.json"),
	}, nil
}

func (s *FileStore) Path() string {
	return s.path
}

func (s *FileStore) Load() (AppConfig, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return defaultConfig(), nil
	}
	if err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{}, err
	}
	if cfg.Clients == nil {
		cfg.Clients = map[string]ClientConnection{}
	}
	if cfg.Version == 0 {
		cfg.Version = 1
	}
	return cfg, nil
}

func (s *FileStore) Save(cfg AppConfig) error {
	if cfg.Clients == nil {
		cfg.Clients = map[string]ClientConnection{}
	}
	if cfg.Version == 0 {
		cfg.Version = 1
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}
