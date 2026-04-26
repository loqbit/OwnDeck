// Package discoverysvc provides MCP / Skill / client discovery to the
// Wails handler layer, fan-out to the configured connector registry.
package discoverysvc

import (
	"context"

	"OwnDeck/internal/connector"
	"OwnDeck/internal/discovery"
)

type Service struct {
	registry *connector.Registry
}

func New(registry *connector.Registry) *Service {
	return &Service{registry: registry}
}

func (s *Service) Clients() []discovery.ClientInfo {
	cs := s.registry.All()
	out := make([]discovery.ClientInfo, 0, len(cs))
	for _, c := range cs {
		info := c.Probe()
		info.Permission = "read"
		out = append(out, info)
	}
	return out
}

func (s *Service) MCPServers(ctx context.Context) ([]discovery.MCPServer, error) {
	return s.registry.MCPForClients(ctx, []string{"codex"})
}

func (s *Service) MCPServersForClients(ctx context.Context, clientIDs []string) ([]discovery.MCPServer, error) {
	return s.registry.MCPForClients(ctx, clientIDs)
}

func (s *Service) SkillsForClients(ctx context.Context, clientIDs []string) ([]discovery.SkillAsset, error) {
	return s.registry.SkillsForClients(ctx, clientIDs)
}
