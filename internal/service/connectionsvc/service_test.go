package connectionsvc

import (
	"sort"
	"testing"

	"OwnDeck/internal/repository/config"
)

// memStore is an in-memory config.Store for testing.
type memStore struct {
	cfg config.AppConfig
}

func newMemStore() *memStore {
	return &memStore{cfg: config.AppConfig{
		Version: 1,
		Clients: map[string]config.ClientConnection{},
	}}
}

func (m *memStore) Path() string                        { return "/test/config.json" }
func (m *memStore) Load() (config.AppConfig, error)     { return m.cfg, nil }
func (m *memStore) Save(cfg config.AppConfig) error      { m.cfg = cfg; return nil }

func TestService_Connect(t *testing.T) {
	store := newMemStore()
	svc := New(store)

	cfg, err := svc.Connect("codex")
	if err != nil {
		t.Fatal(err)
	}

	conn, ok := cfg.Clients["codex"]
	if !ok {
		t.Fatal("codex not found after connect")
	}
	if !conn.Connected {
		t.Error("expected connected = true")
	}
	if conn.Permission != "read" {
		t.Errorf("permission = %q, want %q", conn.Permission, "read")
	}
	if conn.ConnectedAt == "" {
		t.Error("expected non-empty connectedAt")
	}
}

func TestService_Disconnect(t *testing.T) {
	store := newMemStore()
	svc := New(store)

	// Connect first
	if _, err := svc.Connect("codex"); err != nil {
		t.Fatal(err)
	}

	// Then disconnect
	cfg, err := svc.Disconnect("codex")
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := cfg.Clients["codex"]; ok {
		t.Error("codex should have been removed after disconnect")
	}
}

func TestService_ConnectedIDs(t *testing.T) {
	store := newMemStore()
	svc := New(store)

	// Connect two
	if _, err := svc.Connect("codex"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Connect("claude-code"); err != nil {
		t.Fatal(err)
	}

	ids, err := svc.ConnectedIDs()
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(ids)
	if len(ids) != 2 {
		t.Fatalf("expected 2 IDs, got %d", len(ids))
	}
	if ids[0] != "claude-code" || ids[1] != "codex" {
		t.Errorf("ids = %v", ids)
	}
}

func TestService_ConnectedIDs_Empty(t *testing.T) {
	store := newMemStore()
	svc := New(store)

	ids, err := svc.ConnectedIDs()
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty, got %v", ids)
	}
}

func TestService_ConnectOverwrite(t *testing.T) {
	store := newMemStore()
	svc := New(store)

	// Connect twice — should update, not error
	if _, err := svc.Connect("codex"); err != nil {
		t.Fatal(err)
	}
	cfg, err := svc.Connect("codex")
	if err != nil {
		t.Fatal(err)
	}

	// Should still have exactly one entry
	if len(cfg.Clients) != 1 {
		t.Errorf("expected 1 client, got %d", len(cfg.Clients))
	}
}

func TestService_DisconnectNonexistent(t *testing.T) {
	store := newMemStore()
	svc := New(store)

	// Disconnect something that was never connected — should not error
	cfg, err := svc.Disconnect("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Clients) != 0 {
		t.Errorf("expected empty clients, got %d", len(cfg.Clients))
	}
}

func TestService_GetConfig(t *testing.T) {
	store := newMemStore()
	svc := New(store)

	cfg, err := svc.GetConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Version != 1 {
		t.Errorf("version = %d, want 1", cfg.Version)
	}
}

func TestService_Path(t *testing.T) {
	store := newMemStore()
	svc := New(store)
	if svc.Path() != "/test/config.json" {
		t.Errorf("path = %q", svc.Path())
	}
}
