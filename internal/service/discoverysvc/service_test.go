package discoverysvc

import (
	"context"
	"testing"

	"OwnDeck/internal/connector"
	"OwnDeck/internal/discovery"
)

// fakeConnector implements connector.Connector for testing.
type fakeConnector struct {
	id       string
	name     string
	detected bool
	servers  []discovery.MCPServer
	skills   []discovery.SkillAsset
}

func (f *fakeConnector) ID() string   { return f.id }
func (f *fakeConnector) Name() string { return f.name }

func (f *fakeConnector) Probe() discovery.ClientInfo {
	return discovery.ClientInfo{
		ID:       f.id,
		Name:     f.name,
		Detected: f.detected,
		Status:   discovery.ClientStatus(f.detected),
	}
}

func (f *fakeConnector) DiscoverMCP(_ context.Context) ([]discovery.MCPServer, error) {
	return f.servers, nil
}

func (f *fakeConnector) DiscoverSkills(_ context.Context) ([]discovery.SkillAsset, error) {
	return f.skills, nil
}

func TestService_Clients(t *testing.T) {
	c1 := &fakeConnector{id: "a", name: "A", detected: true}
	c2 := &fakeConnector{id: "b", name: "B", detected: false}
	r := connector.NewRegistry(c1, c2)
	svc := New(r)

	clients := svc.Clients()
	if len(clients) != 2 {
		t.Fatalf("expected 2 clients, got %d", len(clients))
	}

	// Both should have read permission
	for _, c := range clients {
		if c.Permission != "read" {
			t.Errorf("client %q permission = %q, want %q", c.ID, c.Permission, "read")
		}
	}

	// Check detection states
	if !clients[0].Detected {
		t.Error("expected client A to be detected")
	}
	if clients[1].Detected {
		t.Error("expected client B to be not detected")
	}
}

func TestService_MCPServersForClients(t *testing.T) {
	c := &fakeConnector{
		id:   "codex",
		name: "Codex",
		servers: []discovery.MCPServer{
			{Name: "test-server", ClientID: "codex"},
		},
	}
	r := connector.NewRegistry(c)
	svc := New(r)

	servers, err := svc.MCPServersForClients(context.Background(), []string{"codex"})
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}
	if servers[0].Name != "test-server" {
		t.Errorf("name = %q", servers[0].Name)
	}
}

func TestService_SkillsForClients(t *testing.T) {
	c := &fakeConnector{
		id:   "codex",
		name: "Codex",
		skills: []discovery.SkillAsset{
			{Name: "git-workflow", ClientID: "codex"},
		},
	}
	r := connector.NewRegistry(c)
	svc := New(r)

	skills, err := svc.SkillsForClients(context.Background(), []string{"codex"})
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

func TestService_MCPServers_DefaultsToCodex(t *testing.T) {
	c := &fakeConnector{
		id:   "codex",
		name: "Codex",
		servers: []discovery.MCPServer{
			{Name: "default", ClientID: "codex"},
		},
	}
	r := connector.NewRegistry(c)
	svc := New(r)

	// MCPServers() uses hardcoded ["codex"] IDs
	servers, err := svc.MCPServers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}
}
