package platform

import "os"

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

func PathExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func ExistingPaths(paths ...string) []string {
	var existing []string
	for _, p := range paths {
		if PathExists(p) {
			existing = append(existing, p)
		}
	}
	return existing
}
