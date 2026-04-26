package discovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestScanAgents_ReturnsKnownAgents(t *testing.T) {
	agents := ScanAgents()
	if len(agents) < 5 {
		t.Fatalf("expected at least 5 built-in agents, got %d", len(agents))
	}

	// Verify all known IDs are present
	ids := map[string]bool{}
	for _, a := range agents {
		ids[a.ID] = true
	}
	for _, expected := range []string{"claude-code", "claude-desktop", "codex", "gemini-cli", "antigravity"} {
		if !ids[expected] {
			t.Errorf("missing expected agent: %s", expected)
		}
	}
}

func TestScanAgents_FieldsPopulated(t *testing.T) {
	agents := ScanAgents()
	for _, a := range agents {
		if a.ID == "" {
			t.Error("agent has empty ID")
		}
		if a.Name == "" {
			t.Errorf("agent %s has empty Name", a.ID)
		}
		if a.ScannedAt == "" {
			t.Errorf("agent %s has empty ScannedAt", a.ID)
		}
	}
}

func TestHasMCPServersKey(t *testing.T) {
	dir := t.TempDir()

	// File with mcpServers key
	withKey := filepath.Join(dir, "with.json")
	data, _ := json.Marshal(map[string]any{
		"mcpServers": map[string]any{
			"test": map[string]any{"command": "echo"},
		},
	})
	_ = os.WriteFile(withKey, data, 0o600)

	// File without mcpServers key
	withoutKey := filepath.Join(dir, "without.json")
	data, _ = json.Marshal(map[string]any{
		"settings": map[string]any{"theme": "dark"},
	})
	_ = os.WriteFile(withoutKey, data, 0o600)

	// Invalid JSON
	invalid := filepath.Join(dir, "invalid.json")
	_ = os.WriteFile(invalid, []byte("not json"), 0o600)

	if !hasMCPServersKey(withKey) {
		t.Error("expected true for file with mcpServers key")
	}
	if hasMCPServersKey(withoutKey) {
		t.Error("expected false for file without mcpServers key")
	}
	if hasMCPServersKey(invalid) {
		t.Error("expected false for invalid JSON")
	}
	if hasMCPServersKey(filepath.Join(dir, "nonexistent.json")) {
		t.Error("expected false for nonexistent file")
	}
}
