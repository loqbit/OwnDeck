# OwnDeck UI References

## Product Shape

OwnDeck should feel like a local control console, not a marketing dashboard.

The closest references are:

- Docker Desktop: local runtime control, left navigation, resource lists, extension area.
- Raycast Preferences / Extensions: installed capabilities, commands, extension-specific settings.
- Tailscale Admin Console: connected devices, services, status, tags, filters.
- Linear Settings / Integrations: clean integration management, restrained visual density.
- Postman / Insomnia / TablePlus: resource navigation, list/detail workflow, action-focused tooling.

## What To Borrow

### Docker Desktop

Use for the overall application shell:

- Left sidebar navigation.
- Main resource pages.
- Clear distinction between runtime objects and extensions.
- Plain, operational UI instead of decorative cards.

Useful pages:

- https://docs.docker.com/desktop/use-desktop/
- https://www.docker.com/products/extensions/

### Raycast

Use for Skills and plugin-like assets:

- Extension list.
- Search-first management.
- Extension details with commands/settings.
- Clear installed vs built-in distinction.

Useful pages:

- https://manual.raycast.com/preferences
- https://manual.raycast.com/windows/settings

### Tailscale

Use for client and service status:

- Device/service list.
- Filters and tags.
- Status-first rows.
- Connection-oriented language.

Useful pages:

- https://tailscale.com/docs/features/access-control/device-management/how-to/filter
- https://tailscale.com/kb/1552/tailscale-services

### Linear

Use for integration settings:

- Quiet integration cards.
- Simple enable/disable states.
- Clear settings hierarchy.

Useful page:

- https://www.saasui.design/pattern/settings/linear

## Recommended OwnDeck Layout

```text
OwnDeck
├── Overview
├── Clients
├── MCP Servers
├── Skills
├── Agents
├── Profiles
└── Settings
```

## Page Direction

### Overview

Purpose: quick health and inventory.

Show:

- Connected clients
- MCP servers
- User skills
- Agents
- Warnings: missing CLI, unreadable config, failed health checks

Avoid:

- Long tables
- Full skill descriptions
- Decorative cards

### Clients

Purpose: connection consent and client detection.

Show:

- Codex
- Claude Code
- Antigravity
- Cursor
- Zed
- VS Code / Cline / Roo later

Each row/card should show:

- Detected / not found
- Connected / disconnected
- Read-only / manage permission
- Executable path
- Config paths
- Connect / Disconnect

### MCP Servers

Purpose: tool server inventory.

Show as table:

- Name
- Client
- Transport
- Command / URL
- Status
- Source config
- Env present

Later:

- Health check
- Enable / disable per profile
- Diff preview before export

### Skills

Purpose: manage useful local skill packages.

Default filter:

- User / third-party skills

Secondary filters:

- Built-in
- Plugin
- All

Show:

- Name
- Client
- Scope
- Description
- Source path

Later:

- Install from GitHub
- Import local folder
- Validate SKILL.md
- Sync to another client

### Agents

Purpose: A2A future surface.

For now:

- Empty state
- Explain Agent Cards will appear here

Later:

- Agent name
- Endpoint
- Capabilities
- Auth state
- Health
- Call / inspect

### Profiles

Purpose: OwnDeck's real product surface.

Show:

- Profile name
- Target clients
- Enabled MCP servers
- Enabled skills
- Future agents

Actions:

- Preview export
- Apply
- Backup
- Rollback

## Visual Principles

- Prefer tables and compact lists for assets.
- Use cards only for connection setup, empty states, and small repeated integration summaries.
- Avoid putting cards inside cards.
- Keep text compact and scannable.
- Use one primary accent color, but do not let the app become a purple dashboard.
- Always show source paths in monospace.
- Separate "detected" from "connected".
- Separate "MCP server" from "Agent" from "Skill".

## Next UI Iteration

Replace the current one-page dashboard with:

```text
Sidebar
  Overview
  Clients
  MCP Servers
  Skills
  Agents

Main
  Header with page title and Refresh
  Page-specific toolbar
  List/table
  Optional detail drawer
```

This should improve both performance and clarity because inactive sections do not need to render large lists.
