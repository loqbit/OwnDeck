// Package discovery defines the wire types returned by the Wails layer.
// Service logic lives in internal/service/discoverysvc; per-client
// implementations live in internal/connector.
package discovery

type MCPServer struct {
	Name       string `json:"name"`
	ClientID   string `json:"clientID"`
	Client     string `json:"client"`
	Transport  string `json:"transport"`
	Command    string `json:"command"`
	Args       string `json:"args"`
	URL        string `json:"url"`
	Env        string `json:"env"`
	Cwd        string `json:"cwd"`
	Status     string `json:"status"`
	Auth       string `json:"auth"`
	SourcePath string `json:"sourcePath"`

	// Origin describes where this server was declared.
	// Examples: "user" (the client's own user-level config) or
	// "plugin:<plugin-id>" (a plugin bundled the server).
	Origin     string `json:"origin"`
	OriginPath string `json:"originPath"`
}

type SkillAsset struct {
	Name        string `json:"name"`
	ClientID    string `json:"clientID"`
	Client      string `json:"client"`
	Description string `json:"description"`
	SourcePath  string `json:"sourcePath"`
	Scope       string `json:"scope"`
}

type ClientInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Detected    bool     `json:"detected"`
	Connected   bool     `json:"connected"`
	Permission  string   `json:"permission"`
	Executable  string   `json:"executable"`
	ConfigPaths []string `json:"configPaths"`
	Status      string   `json:"status"`
}

func ClientStatus(detected bool) string {
	if detected {
		return "detected"
	}
	return "not found"
}
