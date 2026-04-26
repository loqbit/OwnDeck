package platform

import "os/exec"

// LookPath returns the absolute path of name, or "" if not found.
// Used only for Probe (informational); MCP discovery never shells out.
func LookPath(name string) string {
	if path, err := exec.LookPath(name); err == nil {
		return path
	}
	return ""
}
