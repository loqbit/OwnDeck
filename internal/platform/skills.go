package platform

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type SkillFile struct {
	Path        string
	Name        string
	Description string
	Scope       string
}

func DiscoverSkillFiles(roots []string) []SkillFile {
	var skills []SkillFile
	seen := map[string]bool{}

	for _, root := range ExistingPaths(roots...) {
		_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() || entry.Name() != "SKILL.md" || seen[path] {
				return nil
			}
			sf, err := parseSkillFile(path, skillScope(root, path))
			if err != nil {
				return nil
			}
			seen[path] = true
			skills = append(skills, sf)
			return nil
		})
	}
	return skills
}

func parseSkillFile(path string, scope string) (SkillFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SkillFile{}, err
	}

	name := filepath.Base(filepath.Dir(path))
	description := ""
	lines := strings.Split(string(data), "\n")

	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for _, line := range lines[1:] {
			trimmed := strings.TrimSpace(line)
			if trimmed == "---" {
				break
			}

			key, value, ok := strings.Cut(trimmed, ":")
			if !ok {
				continue
			}

			value = strings.Trim(strings.TrimSpace(value), `"'`)
			switch strings.TrimSpace(key) {
			case "name":
				if value != "" {
					name = value
				}
			case "description":
				description = value
			}
		}
	}

	return SkillFile{Path: path, Name: name, Description: description, Scope: scope}, nil
}

func skillScope(root string, path string) string {
	sep := string(filepath.Separator)
	if strings.Contains(path, sep+"plugins"+sep+"cache"+sep) {
		return "plugin"
	}
	if strings.Contains(path, sep+".system"+sep) {
		return "system"
	}
	if strings.HasPrefix(path, root) {
		return "user"
	}
	return "unknown"
}
