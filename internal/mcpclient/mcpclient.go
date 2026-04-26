// Package mcpclient provides a minimal MCP-over-stdio client for
// introspecting MCP servers. It temporarily starts a server process,
// performs the initialize handshake, queries tools/list (and optionally
// prompts/list, resources/list), then tears the connection down.
//
// OwnDeck uses this to populate ToolInfo and health status for each
// discovered MCP server. It never calls tools/call — only read-only
// introspection.
package mcpclient

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const DefaultTimeout = 10 * time.Second

// ToolInfo mirrors discovery.ToolInfo but lives here to avoid an
// import cycle. The caller maps these to discovery types.
type ToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema,omitempty"`
}

// IntrospectResult holds everything gathered from a single
// introspection session with an MCP server.
type IntrospectResult struct {
	Tools         []ToolInfo
	ServerName    string
	ServerVersion string
	Health        string // "healthy" | "degraded" | "error"
	Error         string
	Duration      time.Duration
}

// Introspect starts an MCP server as a subprocess, performs the
// protocol handshake, pulls the tools list, and shuts down.
// The entire operation is bounded by DefaultTimeout.
func Introspect(ctx context.Context, command string, args []string, env map[string]string, cwd string) IntrospectResult {
	return IntrospectWithTimeout(ctx, DefaultTimeout, command, args, env, cwd)
}

// IntrospectWithTimeout is like Introspect but with a configurable
// timeout.
func IntrospectWithTimeout(ctx context.Context, timeout time.Duration, command string, args []string, env map[string]string, cwd string) IntrospectResult {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result := doIntrospect(ctx, command, args, env, cwd)
	result.Duration = time.Since(start)
	return result
}

func doIntrospect(ctx context.Context, command string, args []string, env map[string]string, cwd string) IntrospectResult {
	// Build the subprocess command.
	cmd := exec.CommandContext(ctx, command, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	if len(env) > 0 {
		cmd.Env = buildEnv(env)
	}

	// Create the MCP client and transport.
	client := mcp.NewClient(
		&mcp.Implementation{Name: "owndeck", Version: "0.1.0"},
		nil,
	)
	transport := &mcp.CommandTransport{Command: cmd}

	// Connect — this performs initialize + notifications/initialized.
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return IntrospectResult{
			Health: "error",
			Error:  fmt.Sprintf("connect: %v", err),
		}
	}
	defer session.Close()

	// List tools.
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		return IntrospectResult{
			Health: "degraded",
			Error:  fmt.Sprintf("tools/list: %v", err),
		}
	}

	tools := make([]ToolInfo, 0, len(toolsResult.Tools))
	for _, t := range toolsResult.Tools {
		tools = append(tools, ToolInfo{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		})
	}

	return IntrospectResult{
		Tools:  tools,
		Health: "healthy",
	}
}

// buildEnv merges the provided env vars into the current process
// environment so things like PATH, HOME, etc. are preserved.
func buildEnv(extra map[string]string) []string {
	// Start with os environ would be ideal but exec.CommandContext
	// inherits the parent env by default when cmd.Env is nil.
	// When we set cmd.Env, we need to explicitly include everything.
	// To keep it simple and safe, we build a minimal set by re-
	// using the system env.
	var env []string
	// Import key system vars
	for _, key := range []string{"PATH", "HOME", "USER", "SHELL", "TERM", "LANG", "TMPDIR"} {
		if val := os.Getenv(key); val != "" {
			env = append(env, key+"="+val)
		}
	}
	for k, v := range extra {
		env = append(env, k+"="+v)
	}
	return env
}

