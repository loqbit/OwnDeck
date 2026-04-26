package connector

import (
	"context"
	"errors"
	"testing"

	"OwnDeck/internal/discovery"
)

// fakeConnector implements Connector for testing the Registry.
type fakeConnector struct {
	id       string
	name     string
	info     discovery.ClientInfo
	servers  []discovery.MCPServer
	skills   []discovery.SkillAsset
	mcpErr   error
	skillErr error
}

func (f *fakeConnector) ID() string   { return f.id }
func (f *fakeConnector) Name() string { return f.name }

func (f *fakeConnector) Probe() discovery.ClientInfo {
	return f.info
}

func (f *fakeConnector) DiscoverMCP(_ context.Context) ([]discovery.MCPServer, error) {
	return f.servers, f.mcpErr
}

func (f *fakeConnector) DiscoverSkills(_ context.Context) ([]discovery.SkillAsset, error) {
	return f.skills, f.skillErr
}

func TestRegistry_All(t *testing.T) {
	c1 := &fakeConnector{id: "a", name: "A"}
	c2 := &fakeConnector{id: "b", name: "B"}
	r := NewRegistry(c1, c2)

	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 connectors, got %d", len(all))
	}
}

func TestRegistry_Get(t *testing.T) {
	c := &fakeConnector{id: "test", name: "Test"}
	r := NewRegistry(c)

	found, ok := r.Get("test")
	if !ok || found.ID() != "test" {
		t.Error("failed to find connector 'test'")
	}

	_, ok = r.Get("missing")
	if ok {
		t.Error("found connector that should not exist")
	}
}

func TestRegistry_MCPForClients(t *testing.T) {
	c1 := &fakeConnector{
		id:   "codex",
		name: "Codex",
		servers: []discovery.MCPServer{
			{Name: "fs", ClientID: "codex"},
			{Name: "github", ClientID: "codex"},
		},
	}
	c2 := &fakeConnector{
		id:   "claude",
		name: "Claude",
		servers: []discovery.MCPServer{
			{Name: "web", ClientID: "claude"},
		},
	}
	r := NewRegistry(c1, c2)

	// Request both
	servers, err := r.MCPForClients(context.Background(), []string{"codex", "claude"})
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 3 {
		t.Fatalf("expected 3 servers, got %d", len(servers))
	}

	// Request one
	servers, err = r.MCPForClients(context.Background(), []string{"claude"})
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}
}

func TestRegistry_MCPForClients_Deduplicates(t *testing.T) {
	c := &fakeConnector{
		id:      "test",
		name:    "Test",
		servers: []discovery.MCPServer{{Name: "s"}},
	}
	r := NewRegistry(c)

	servers, err := r.MCPForClients(context.Background(), []string{"test", "test", "test"})
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected 1 server (deduplicated), got %d", len(servers))
	}
}

func TestRegistry_MCPForClients_UnsupportedClient(t *testing.T) {
	r := NewRegistry() // empty registry

	// All unsupported → error
	_, err := r.MCPForClients(context.Background(), []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for unsupported client")
	}
}

func TestRegistry_MCPForClients_PartialFailure(t *testing.T) {
	good := &fakeConnector{
		id:      "good",
		name:    "Good",
		servers: []discovery.MCPServer{{Name: "ok"}},
	}
	bad := &fakeConnector{
		id:     "bad",
		name:   "Bad",
		mcpErr: errors.New("boom"),
	}
	r := NewRegistry(good, bad)

	// Should return the good results without error
	servers, err := r.MCPForClients(context.Background(), []string{"good", "bad"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected 1 server from good connector, got %d", len(servers))
	}
}

func TestRegistry_MCPForClients_AllFail(t *testing.T) {
	bad := &fakeConnector{
		id:     "bad",
		name:   "Bad",
		mcpErr: errors.New("boom"),
	}
	r := NewRegistry(bad)

	servers, err := r.MCPForClients(context.Background(), []string{"bad"})
	if err == nil {
		t.Fatal("expected error when all connectors fail")
	}
	if len(servers) != 0 {
		t.Errorf("expected 0 servers, got %d", len(servers))
	}
}

func TestRegistry_SkillsForClients(t *testing.T) {
	c := &fakeConnector{
		id:   "codex",
		name: "Codex",
		skills: []discovery.SkillAsset{
			{Name: "git-workflow", ClientID: "codex"},
		},
	}
	r := NewRegistry(c)

	skills, err := r.SkillsForClients(context.Background(), []string{"codex"})
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Name != "git-workflow" {
		t.Errorf("name = %q", skills[0].Name)
	}
}

func TestRegistry_EmptyResults_NonNilSlice(t *testing.T) {
	c := &fakeConnector{
		id:      "empty",
		name:    "Empty",
		servers: nil, // DiscoverMCP returns nil
	}
	r := NewRegistry(c)

	servers, err := r.MCPForClients(context.Background(), []string{"empty"})
	if err != nil {
		t.Fatal(err)
	}
	// aggregate() initialises with make([]T, 0), so even with nil
	// from the connector, the result must be non-nil ([] not null in JSON).
	if servers == nil {
		t.Error("expected non-nil (empty) slice, got nil")
	}
}
