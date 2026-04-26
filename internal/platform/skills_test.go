package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverSkillFiles_WithFrontmatter(t *testing.T) {
	root := t.TempDir()
	skillDir := filepath.Join(root, "my-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := `---
name: "Test Skill"
description: "A skill for testing"
---
# Test Skill
This is a test skill.
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	skills := DiscoverSkillFiles([]string{root})
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}

	s := skills[0]
	if s.Name != "Test Skill" {
		t.Errorf("name = %q, want %q", s.Name, "Test Skill")
	}
	if s.Description != "A skill for testing" {
		t.Errorf("description = %q, want %q", s.Description, "A skill for testing")
	}
	if s.Path != filepath.Join(skillDir, "SKILL.md") {
		t.Errorf("path = %q", s.Path)
	}
}

func TestDiscoverSkillFiles_NoFrontmatter(t *testing.T) {
	root := t.TempDir()
	skillDir := filepath.Join(root, "plain-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := `# Plain Skill
No frontmatter here.
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	skills := DiscoverSkillFiles([]string{root})
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}

	// Falls back to parent directory name
	if skills[0].Name != "plain-skill" {
		t.Errorf("name = %q, want %q", skills[0].Name, "plain-skill")
	}
	if skills[0].Description != "" {
		t.Errorf("description = %q, want empty", skills[0].Description)
	}
}

func TestDiscoverSkillFiles_NestedSkills(t *testing.T) {
	root := t.TempDir()
	// Create two skills at different nesting levels
	skill1 := filepath.Join(root, "skills", "level1")
	skill2 := filepath.Join(root, "skills", "nested", "level2")
	for _, dir := range []string{skill1, skill2} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# Skill\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	skills := DiscoverSkillFiles([]string{root})
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
}

func TestDiscoverSkillFiles_EmptyRoots(t *testing.T) {
	skills := DiscoverSkillFiles(nil)
	if skills != nil {
		t.Errorf("expected nil, got %v", skills)
	}
}

func TestDiscoverSkillFiles_NonexistentRoot(t *testing.T) {
	skills := DiscoverSkillFiles([]string{"/nonexistent/path"})
	if skills != nil {
		t.Errorf("expected nil, got %v", skills)
	}
}

func TestDiscoverSkillFiles_DeduplicatesAcrossRoots(t *testing.T) {
	root := t.TempDir()
	skillDir := filepath.Join(root, "skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Pass the same root twice
	skills := DiscoverSkillFiles([]string{root, root})
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill (deduplicated), got %d", len(skills))
	}
}

func TestSkillScope_Plugin(t *testing.T) {
	scope := skillScope("/root", "/root/plugins/cache/marketplace/plugin/v1/SKILL.md")
	if scope != "plugin" {
		t.Errorf("scope = %q, want %q", scope, "plugin")
	}
}

func TestSkillScope_System(t *testing.T) {
	scope := skillScope("/root", "/root/.system/builtins/SKILL.md")
	if scope != "system" {
		t.Errorf("scope = %q, want %q", scope, "system")
	}
}

func TestSkillScope_User(t *testing.T) {
	scope := skillScope("/root/skills", "/root/skills/my-skill/SKILL.md")
	if scope != "user" {
		t.Errorf("scope = %q, want %q", scope, "user")
	}
}

func TestSkillScope_Unknown(t *testing.T) {
	scope := skillScope("/root", "/other/path/SKILL.md")
	if scope != "unknown" {
		t.Errorf("scope = %q, want %q", scope, "unknown")
	}
}
