package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHomeDir(t *testing.T) {
	home := HomeDir()
	if home == "" {
		t.Fatal("HomeDir returned empty string")
	}
	// Should be an absolute path
	if !filepath.IsAbs(home) {
		t.Errorf("HomeDir returned non-absolute path: %q", home)
	}
}

func TestPathExists_ExistingPath(t *testing.T) {
	dir := t.TempDir()
	if !PathExists(dir) {
		t.Errorf("PathExists(%q) = false, want true", dir)
	}
}

func TestPathExists_NonexistentPath(t *testing.T) {
	if PathExists("/definitely/does/not/exist/12345") {
		t.Error("PathExists returned true for nonexistent path")
	}
}

func TestPathExists_EmptyString(t *testing.T) {
	if PathExists("") {
		t.Error("PathExists returned true for empty string")
	}
}

func TestExistingPaths_Mixed(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "exists.txt")
	if err := os.WriteFile(file, []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := ExistingPaths(file, "/no/such/path", dir)
	if len(result) != 2 {
		t.Fatalf("expected 2 existing paths, got %d: %v", len(result), result)
	}
	if result[0] != file {
		t.Errorf("result[0] = %q, want %q", result[0], file)
	}
	if result[1] != dir {
		t.Errorf("result[1] = %q, want %q", result[1], dir)
	}
}

func TestExistingPaths_NoneExist(t *testing.T) {
	result := ExistingPaths("/a/b/c", "/d/e/f")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestExistingPaths_Empty(t *testing.T) {
	result := ExistingPaths()
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}
