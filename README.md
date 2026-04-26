# OwnDeck

Local-first AI Agent Control Plane.

OwnDeck is the only local-first console that actually speaks the MCP protocol — it manages MCP servers (down to the individual tool), Skills, and (later) A2A agents in one place, then compiles them into Claude / Codex / Cursor / Zed configuration on demand. Existing tools (MCP Linker, mcpm, Vibe Manager, MCPX, MyMCP) only synchronize JSON between clients; OwnDeck introspects each server, knows what tools it exposes, and treats Skills and Agents as first-class assets alongside MCP.

## Development

Requirements:

- Go 1.21+
- Node / npm
- Wails v2

Install Wails:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails doctor
```

Run the desktop app in development mode:

```bash
wails dev
```

Build a local app bundle:

```bash
wails build
```

Run tests:

```bash
go test ./...
```

## Current Scope (Stage 1 Complete)

- Wails + Vue 3 + TypeScript + shadcn-vue app shell with sidebar navigation
- Layered Go backend (`connector / service / repository / platform / discovery`) with explicit dependency injection in `main.go`
- **Direct config-file reading** — no CLI subprocess scraping anywhere:
  - **Codex**: reads `~/.codex/config.toml` (TOML) plus plugin-bundled `.mcp.json` files from `~/.codex/plugins/cache/`; each server is tagged with `Origin` ("user" vs "plugin:\<id\>")
  - **Claude Code**: reads `~/.claude/settings.json`, `~/.claude.json`, and project-level `.mcp.json`
  - **Claude Desktop**: reads `~/Library/Application Support/Claude/claude_desktop_config.json`
  - **Antigravity**: reads `~/.gemini/antigravity/mcp_config.json` and other JSON config candidates
- Read-only Skill discovery from Codex and Antigravity `SKILL.md` packages
- Skill filtering: user-managed vs bundled system/plugin skills
- Multi-page UI: Overview, Clients, MCP Servers, Skills, Agents (placeholder), Profiles (placeholder), Settings
- Persisted client connection consent in OwnDeck's local config
- i18n (English / Chinese)
- Auto-refreshing data polling with silent background updates

Not yet implemented (see [OwnDeck.md](OwnDeck.md) for the full roadmap):

- Native MCP-over-stdio client for live tool / prompt / resource introspection (Stage 2)
- Health checks and health status display (Stage 2)
- SQLite-backed Profile & Asset model (Stage 3)
- Multi-client config compilation with diff preview, backup, and rollback (Stage 4)
- Tool-level enable/disable, cross-client version drift detection (Stage 2–3)
- Cursor connector (standard JSON format, deferred until needed)
- A2A agent registry (Stage 5+)

OwnDeck stores its own local app configuration under the user's config directory:

```text
~/Library/Application Support/OwnDeck/config.json
```

This file records client connection consent only. OwnDeck does not write external client MCP config files yet.

## Client Discovery Model

OwnDeck cannot reliably discover every AI client generically because local AI tools do not share a common registration standard for MCP config paths, CLI commands, scopes, or project files.

Instead, OwnDeck uses a connector registry. Each connector owns its own probe and discovery logic:

- probe whether the client exists
- report executable and config paths
- read MCP servers after the user connects the client
- later, export config through the safest client-specific mechanism

This keeps discovery explicit and makes new clients like Cursor, Zed, and VS Code additive instead of special-casing them in the UI.

## Notable Implementation Details

- **Codex plugin-bundled MCP servers**: OwnDeck walks enabled plugins' `.codex-plugin/plugin.json` manifests and follows their `mcpServers` path references to discover bundled servers. This is something the Codex desktop UI itself currently fails to do correctly (see [openai/codex#17360](https://github.com/openai/codex/issues/17360)).
- **Origin tagging**: Every discovered MCP server carries an `Origin` field (e.g., `"user"`, `"plugin:github@openai-curated"`, `"project"`) so the UI can group and filter them by provenance.
- **No CLI dependency**: All discovery is done by reading configuration files directly (TOML for Codex, JSON for everything else). The CLI does not need to be installed for OwnDeck to work.
