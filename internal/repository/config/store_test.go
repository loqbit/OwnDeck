package config

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestFileStore_LoadSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := &FileStore{path: filepath.Join(dir, "config.json")}

	// Initial load should return default config
	cfg, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Version != 1 {
		t.Errorf("version = %d, want 1", cfg.Version)
	}
	if len(cfg.Clients) != 0 {
		t.Errorf("expected empty clients, got %d", len(cfg.Clients))
	}

	// Save with data
	cfg.Clients["codex"] = ClientConnection{
		Connected:  true,
		Permission: "read",
	}
	if err := store.Save(cfg); err != nil {
		t.Fatal(err)
	}

	// Reload
	cfg2, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg2.Version != 1 {
		t.Errorf("version = %d, want 1", cfg2.Version)
	}
	conn, ok := cfg2.Clients["codex"]
	if !ok {
		t.Fatal("codex not found in loaded config")
	}
	if !conn.Connected {
		t.Error("expected codex connected = true")
	}
	if conn.Permission != "read" {
		t.Errorf("permission = %q, want %q", conn.Permission, "read")
	}
}

func TestFileStore_Load_MissingFile(t *testing.T) {
	store := &FileStore{path: filepath.Join(t.TempDir(), "nonexistent", "config.json")}
	cfg, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	// Should return default
	if cfg.Version != 1 {
		t.Errorf("version = %d, want 1", cfg.Version)
	}
}

func TestFileStore_Save_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "config.json")
	store := &FileStore{path: path}

	cfg := AppConfig{Version: 1, Clients: map[string]ClientConnection{}}
	if err := store.Save(cfg); err != nil {
		t.Fatal(err)
	}

	// Verify file was created
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestFileStore_Load_NilClientsInitialized(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	// Write JSON with null clients
	if err := os.WriteFile(path, []byte(`{"version": 1, "clients": null}`), 0o644); err != nil {
		t.Fatal(err)
	}

	store := &FileStore{path: path}
	cfg, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Clients == nil {
		t.Error("expected non-nil Clients map")
	}
}

func TestFileStore_Load_ZeroVersionFixed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(`{"clients": {}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	store := &FileStore{path: path}
	cfg, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Version != 1 {
		t.Errorf("version = %d, want 1 (should be auto-fixed)", cfg.Version)
	}
}

func TestFileStore_Load_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("{bad json}"), 0o644); err != nil {
		t.Fatal(err)
	}

	store := &FileStore{path: path}
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFileStore_Path(t *testing.T) {
	store := &FileStore{path: "/test/path/config.json"}
	if store.Path() != "/test/path/config.json" {
		t.Errorf("path = %q", store.Path())
	}
}

func TestConnectedClientIDs(t *testing.T) {
	cfg := AppConfig{
		Clients: map[string]ClientConnection{
			"codex":          {Connected: true},
			"claude-code":    {Connected: true},
			"claude-desktop": {Connected: false},
			"antigravity":    {Connected: true},
		},
	}

	ids := ConnectedClientIDs(cfg)
	sort.Strings(ids)
	if len(ids) != 3 {
		t.Fatalf("expected 3 connected, got %d: %v", len(ids), ids)
	}
	expected := []string{"antigravity", "claude-code", "codex"}
	for i, id := range expected {
		if ids[i] != id {
			t.Errorf("ids[%d] = %q, want %q", i, ids[i], id)
		}
	}
}

func TestConnectedClientIDs_Empty(t *testing.T) {
	cfg := AppConfig{
		Clients: map[string]ClientConnection{},
	}
	ids := ConnectedClientIDs(cfg)
	if len(ids) != 0 {
		t.Errorf("expected empty, got %v", ids)
	}
}

func TestConnectedClientIDs_AllDisconnected(t *testing.T) {
	cfg := AppConfig{
		Clients: map[string]ClientConnection{
			"a": {Connected: false},
			"b": {Connected: false},
		},
	}
	ids := ConnectedClientIDs(cfg)
	if len(ids) != 0 {
		t.Errorf("expected empty, got %v", ids)
	}
}
