// Package connector defines the Connector interface that adapts a local
// AI client (Codex, Claude Code, Antigravity, ...) into OwnDeck's
// discovery model, plus a Registry for routing to them by ID.
package connector

import (
	"context"
	"errors"
	"strings"

	"OwnDeck/internal/discovery"
)

type Connector interface {
	ID() string
	Name() string
	Probe() discovery.ClientInfo
	DiscoverMCP(ctx context.Context) ([]discovery.MCPServer, error)
	DiscoverSkills(ctx context.Context) ([]discovery.SkillAsset, error)
}

type Registry struct {
	connectors []Connector
	byID       map[string]Connector
}

func NewRegistry(connectors ...Connector) *Registry {
	byID := make(map[string]Connector, len(connectors))
	for _, c := range connectors {
		byID[c.ID()] = c
	}
	return &Registry{connectors: connectors, byID: byID}
}

func (r *Registry) All() []Connector {
	return r.connectors
}

func (r *Registry) Get(id string) (Connector, bool) {
	c, ok := r.byID[id]
	return c, ok
}

func (r *Registry) MCPForClients(ctx context.Context, clientIDs []string) ([]discovery.MCPServer, error) {
	return aggregate(r, clientIDs, func(c Connector) ([]discovery.MCPServer, error) {
		return c.DiscoverMCP(ctx)
	})
}

func (r *Registry) SkillsForClients(ctx context.Context, clientIDs []string) ([]discovery.SkillAsset, error) {
	return aggregate(r, clientIDs, func(c Connector) ([]discovery.SkillAsset, error) {
		return c.DiscoverSkills(ctx)
	})
}

func aggregate[T any](r *Registry, clientIDs []string, fn func(Connector) ([]T, error)) ([]T, error) {
	seen := map[string]bool{}
	// Always start with a non-nil slice so callers (and the JSON
	// layer crossing Wails into the frontend) see [] instead of
	// null when no items match.
	results := make([]T, 0)
	var errs []string

	for _, id := range clientIDs {
		if seen[id] {
			continue
		}
		seen[id] = true

		c, ok := r.Get(id)
		if !ok {
			errs = append(errs, id+": unsupported client")
			continue
		}
		items, err := fn(c)
		if err != nil {
			errs = append(errs, c.Name()+": "+err.Error())
			continue
		}
		results = append(results, items...)
	}

	if len(results) == 0 && len(errs) > 0 {
		return results, errors.New(strings.Join(errs, "; "))
	}
	return results, nil
}
