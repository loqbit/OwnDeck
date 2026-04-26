// Package discoverysvc provides MCP / Skill / client discovery to the
// Wails handler layer, fan-out to the configured connector registry.
package discoverysvc

import (
	"context"
	"strings"
	"time"

	"OwnDeck/internal/connector"
	"OwnDeck/internal/discovery"
	"OwnDeck/internal/mcpclient"
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

// IntrospectServer performs MCP protocol introspection on a single
// stdio server. It starts the server process, runs initialize +
// tools/list, then shuts down. The returned MCPServer has its
// Tools, ToolCount, HealthStatus, and IntrospectedAt fields filled.
func (s *Service) IntrospectServer(ctx context.Context, server discovery.MCPServer) discovery.MCPServer {
	if server.Transport != "stdio" || server.Command == "" {
		server.HealthStatus = "error"
		server.HealthMessage = "only stdio servers can be introspected"
		server.IntrospectedAt = time.Now().UTC().Format(time.RFC3339)
		return server
	}

	var args []string
	if server.Args != "" {
		args = strings.Fields(server.Args)
	}

	result := mcpclient.Introspect(ctx, server.Command, args, nil, server.Cwd)

	server.HealthStatus = result.Health
	server.HealthMessage = result.Error
	server.IntrospectedAt = time.Now().UTC().Format(time.RFC3339)

	tools := make([]discovery.ToolInfo, 0, len(result.Tools))
	for _, t := range result.Tools {
		tools = append(tools, discovery.ToolInfo{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		})
	}
	server.Tools = tools
	server.ToolCount = len(tools)

	if result.Health == "healthy" {
		server.Status = "healthy"
	} else {
		server.Status = result.Health
	}

	return server
}

