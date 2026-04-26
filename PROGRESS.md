# Progress

## 2026-04-25

- Reframed OwnDeck as a local AI Agent Control Plane.
- Added the Wails v2 + Vue + TypeScript application shell.
- Added read-only Codex MCP discovery through the Codex CLI.
- Replaced the default Wails sample UI with an OwnDeck MCP discovery table.
- Verified `go test ./...` and `wails build`.
- Added the first client connection flow with read-only Connect / Disconnect actions.
- Added Claude Code detection through `~/.claude/settings.json`, `~/.claude.json`, and project `.mcp.json` candidates.
- Added Antigravity detection through `/Applications/Antigravity.app`, `~/.gemini/antigravity/mcp_config.json`, and project MCP config candidates.
- Updated the UI to show Codex, Claude Code, and Antigravity as separate client cards before loading unified MCP servers.
- Fixed the connection UX so a client can show as connected even when its readable config has zero MCP servers.
- Updated Claude Code discovery to prefer the installed `claude` CLI and fall back to JSON config candidates when the CLI is unavailable.
- Refactored discovery into a connector registry so supported clients provide their own probe and MCP discovery behavior.
- Added a local OwnDeck config store for persisted read-only client connection consent.
- Updated the UI to load saved connections on startup and call backend Connect / Disconnect methods.
- Split the dashboard into Clients, MCP Servers, and Agents sections so tools and agent assets are visually separate.
- Added read-only Skill discovery for connected Codex and Antigravity clients and a dedicated Skills panel in the dashboard.
- Added Skills filters so user-managed skills are emphasized while bundled system/plugin skills are still inspectable.

## 2026-04-26

- Refactored the Go backend into a layered structure: `internal/connector/{codex,claudecode,antigravity}`, `internal/service/{discoverysvc,connectionsvc}`, `internal/repository/config`, `internal/discovery` (types only), `internal/platform`. Added a `Connector` interface + `Registry` and a `config.Store` interface so the App handler is now a thin Wails-bound layer with explicit dependency injection in `main.go`.
- Surveyed the competitive landscape (MCP Linker, mcpm.sh, Vibe Manager, MCPX, MyMCP, MCP Gateway Registry). The "unified MCP profile compiler" pitch is already commoditized, so OwnDeck repivots toward "do all four layers" — MCP protocol introspection, Skills, A2A, and tool-level granularity — as the differentiation.
- Rewrote `OwnDeck.md`: new short-term positioning, explicit Layer A / Layer B architecture, MCP-as-Client introspection chapter, no-CLI-scraping principle, refreshed competitor table, revised five-stage development plan (direct config read → MCP introspection → Profile/Asset → exporters → polish), and updated launch tagline.
- Decision: drop `codex mcp list` / `claude mcp list` subprocess scraping entirely; every connector will read configuration files directly (TOML for Codex, JSON for the rest) and OwnDeck will speak MCP itself for live capability data.
- Replaced Codex's CLI scraping with a direct `~/.codex/config.toml` reader plus a walk over enabled plugins' bundled `.mcp.json` files; tagged each server with an `Origin` field (`user` vs `plugin:<id>`) so the UI can later distinguish the two — something Codex's own desktop UI fails to do (openai/codex#17360).
- Replaced Claude Code's CLI scraping with direct JSON reads of `~/.claude/settings.json`, `~/.claude.json`, and project-local `.mcp.json`.
- Added a Claude Desktop connector reading `~/Library/Application Support/Claude/claude_desktop_config.json`; trimmed the platform layer to only `LookPath` (no more shelling out for discovery).
- Fixed a Vue patch crash on the MCP servers page caused by a `v-for` `<TableRow>` and a `v-if` `<TableEmpty>` (also a `<TableRow>`) sitting as sibling children of `<TableBody>` — switched to `v-if / v-else` mutual exclusion and broadened the row key to include `sourcePath`. Pattern logged in `docs/frontend-pitfalls.md` for future reference.
