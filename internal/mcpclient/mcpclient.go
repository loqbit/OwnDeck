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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const DefaultTimeout = 10 * time.Second

// Health classifications for an introspection attempt.
//
// Distinction matters because OwnDeck inspects servers without ever
// calling tools — a server that responds with a JSON-RPC error has
// proven its protocol layer works, even if it refuses our request
// (e.g. the host app holds the OS permissions needed to actually
// run a tool). Only true silence (timeout, crash, no protocol
// response) means the server is broken from OwnDeck's perspective.
const (
	HealthHealthy     = "healthy"     // initialize OK + tools/list OK
	HealthReachable   = "reachable"   // server spoke MCP but returned an error
	HealthUnreachable = "unreachable" // process never responded / timed out / crashed
)

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
	Health        string // one of HealthHealthy / HealthReachable / HealthUnreachable
	Error         string
	Duration      time.Duration
}

// classify maps an error from Connect or ListTools into a health
// bucket. A JSON-RPC error means the server responded with a valid
// protocol message — it is alive ("reachable"). Timeout, context
// cancel, or any other transport-layer failure means we never heard
// back ("unreachable").
func classify(err error) (health, msg string) {
	var rpcErr *jsonrpc.Error
	if errors.As(err, &rpcErr) {
		return HealthReachable, err.Error()
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return HealthUnreachable, "timeout: " + err.Error()
	}
	return HealthUnreachable, err.Error()
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
		health, msg := classify(err)
		return IntrospectResult{
			Health: health,
			Error:  fmt.Sprintf("connect: %s", msg),
		}
	}
	defer session.Close()

	// List tools.
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		// Init succeeded so the server is at minimum reachable;
		// classify() may still return unreachable if the stream
		// died mid-session (timeout / transport break).
		health, msg := classify(err)
		return IntrospectResult{
			Health: health,
			Error:  fmt.Sprintf("tools/list: %s", msg),
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
		Health: HealthHealthy,
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

