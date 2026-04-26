package platform

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestParseMCPServersFile_StdioServer(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mcp.json")
	content := `{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "env": {"HOME": "/Users/test"}
    }
  }
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	servers, err := ParseMCPServersFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}

	s := servers[0]
	if s.Name != "filesystem" {
		t.Errorf("name = %q, want %q", s.Name, "filesystem")
	}
	if s.Command != "npx" {
		t.Errorf("command = %q, want %q", s.Command, "npx")
	}
	if s.Args != "-y @modelcontextprotocol/server-filesystem /tmp" {
		t.Errorf("args = %q", s.Args)
	}
	if s.Transport != "stdio" {
		t.Errorf("transport = %q, want %q", s.Transport, "stdio")
	}
	if s.Env != "set" {
		t.Errorf("env = %q, want %q", s.Env, "set")
	}
	if s.Source != path {
		t.Errorf("source = %q, want %q", s.Source, path)
	}
}

func TestParseMCPServersFile_HTTPServer(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mcp.json")
	content := `{
  "mcpServers": {
    "remote": {
      "url": "https://example.com/mcp"
    }
  }
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	servers, err := ParseMCPServersFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}

	s := servers[0]
	if s.URL != "https://example.com/mcp" {
		t.Errorf("url = %q", s.URL)
	}
	if s.Transport != "http" {
		t.Errorf("transport = %q, want %q", s.Transport, "http")
	}
}

func TestParseMCPServersFile_MultipleServers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mcp.json")
	content := `{
  "mcpServers": {
    "alpha": {"command": "alpha-cmd"},
    "beta": {"url": "https://beta.example.com"},
    "gamma": {"command": "gamma-cmd", "args": ["--verbose"]}
  }
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	servers, err := ParseMCPServersFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 3 {
		t.Fatalf("expected 3 servers, got %d", len(servers))
	}

	// Sort by name for deterministic assertions
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Name < servers[j].Name
	})

	if servers[0].Name != "alpha" || servers[1].Name != "beta" || servers[2].Name != "gamma" {
		t.Errorf("unexpected server names: %v, %v, %v", servers[0].Name, servers[1].Name, servers[2].Name)
	}
	if servers[1].Transport != "http" {
		t.Errorf("beta transport = %q, want http", servers[1].Transport)
	}
}

func TestParseMCPServersFile_NoMCPServersKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"preferences": {"theme": "dark"}}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	servers, err := ParseMCPServersFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if servers != nil {
		t.Errorf("expected nil, got %v", servers)
	}
}

func TestParseMCPServersFile_EmptyMCPServers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"mcpServers": {}}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	servers, err := ParseMCPServersFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 0 {
		t.Errorf("expected 0 servers, got %d", len(servers))
	}
}

func TestParseMCPServersFile_FileNotFound(t *testing.T) {
	_, err := ParseMCPServersFile("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseMCPServersFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{not valid json}"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseMCPServersFile(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseMCPServersFile_NoEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mcp.json")
	content := `{"mcpServers": {"test": {"command": "test-cmd"}}}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	servers, err := ParseMCPServersFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if servers[0].Env != "-" {
		t.Errorf("env = %q, want %q", servers[0].Env, "-")
	}
}
