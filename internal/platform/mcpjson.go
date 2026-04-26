package platform

import (
	"encoding/json"
	"os"
	"strings"
)

type MCPServerJSON struct {
	Name      string
	Command   string
	Args      string
	URL       string
	Cwd       string
	Env       string
	Transport string
	Source    string
}

type rawMCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	URL     string            `json:"url"`
	Cwd     string            `json:"cwd"`
	Env     map[string]string `json:"env"`
}

func ParseMCPServersFile(path string) ([]MCPServerJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	rawServers, ok := root["mcpServers"]
	if !ok {
		return nil, nil
	}

	var configs map[string]rawMCPServer
	if err := json.Unmarshal(rawServers, &configs); err != nil {
		return nil, err
	}

	out := make([]MCPServerJSON, 0, len(configs))
	for name, cfg := range configs {
		out = append(out, MCPServerJSON{
			Name:      name,
			Command:   cfg.Command,
			Args:      strings.Join(cfg.Args, " "),
			URL:       cfg.URL,
			Cwd:       cfg.Cwd,
			Env:       envSummary(cfg.Env),
			Transport: inferTransport(cfg.URL, cfg.Command),
			Source:    path,
		})
	}
	return out, nil
}

func envSummary(env map[string]string) string {
	if len(env) == 0 {
		return "-"
	}
	return "set"
}

func inferTransport(url, command string) string {
	if url != "" {
		return "http"
	}
	if command != "" {
		return "stdio"
	}
	return "unknown"
}
