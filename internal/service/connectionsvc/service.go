// Package connectionsvc manages user-granted client connection consent
// on top of a config.Store.
package connectionsvc

import (
	"time"

	"OwnDeck/internal/repository/config"
)

type Service struct {
	store config.Store
}

func New(store config.Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetConfig() (config.AppConfig, error) {
	return s.store.Load()
}

func (s *Service) Path() string {
	return s.store.Path()
}

func (s *Service) ConnectedIDs() ([]string, error) {
	cfg, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	return config.ConnectedClientIDs(cfg), nil
}

func (s *Service) Connect(clientID string) (config.AppConfig, error) {
	cfg, err := s.store.Load()
	if err != nil {
		return config.AppConfig{}, err
	}
	cfg.Clients[clientID] = config.ClientConnection{
		Connected:   true,
		Permission:  "read",
		ConnectedAt: time.Now().UTC().Format(time.RFC3339),
	}
	return cfg, s.store.Save(cfg)
}

func (s *Service) Disconnect(clientID string) (config.AppConfig, error) {
	cfg, err := s.store.Load()
	if err != nil {
		return config.AppConfig{}, err
	}
	delete(cfg.Clients, clientID)
	return cfg, s.store.Save(cfg)
}
