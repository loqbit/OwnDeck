<h1 align="center">OwnDeck</h1>

<p align="center">
  <strong>The Local-First AI Agent Control Plane</strong>
</p>

<p align="center">
  <a href="https://wails.io/"><img src="https://img.shields.io/badge/Wails-v2-blue?style=flat-square&logo=go" alt="Wails"></a>
  <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/Vue.js-3.0-4FC08D?style=flat-square&logo=vuedotjs" alt="Vue 3"></a>
  <a href="https://www.typescriptlang.org/"><img src="https://img.shields.io/badge/TypeScript-Ready-3178C6?style=flat-square&logo=typescript" alt="TypeScript"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square" alt="License"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#getting-started">Getting Started</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#contributing">Contributing</a>
</p>

<p align="center">
  <em><a href="README.md">English</a> • <a href="README_zh.md">简体中文</a></em>
</p>

---

OwnDeck is the ultimate local-first console for the Model Context Protocol (MCP). It manages your MCP servers, discovers AI tools, and centralizes your local agent ecosystem in one unified dashboard.

While existing tools merely synchronize static JSON configurations, OwnDeck introduces **protocol-level introspection**. It actively communicates with your servers over `stdio` to map out exact tool capabilities, prompts, and resources.

## Features

- **Dynamic Agent Discovery**: Automatically scans your filesystem to detect installed AI agents (Claude, Codex, Gemini CLI, Cursor) and their extensions.
- **Protocol-Level Introspection**: Directly communicates with MCP-over-stdio servers to live-query tools, prompts, and resources. 
- **Direct Config Parsing**: No fragile CLI subprocess scraping. OwnDeck reads configurations natively (JSON/TOML) with zero external dependencies.
- **Centralized Dashboard**: View all your MCP servers, their connection health, and available tools across your entire machine in a modern interface.
- **Auto-Refresh Architecture**: Built with a reactive stack (Vue 3 + Wails) featuring silent background polling and smooth state management.

## Getting Started

### Prerequisites

Ensure you have the following installed on your system:
- [Go](https://golang.org/doc/install) 1.21 or higher
- [Node.js](https://nodejs.org/) & npm
- [Wails v2](https://wails.io/docs/gettingstarted/installation)

### Installation

1. Install the Wails CLI and verify your environment:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   wails doctor
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/loqbit/OwnDeck.git
   cd OwnDeck
   ```

3. Start the development server:
   ```bash
   wails dev
   ```

4. Build for production:
   ```bash
   wails build
   ```

## Architecture

OwnDeck is built with extensibility in mind, separating presentation from discovery logic.

- **Frontend**: Vue 3 + TypeScript + Tailwind CSS + shadcn-vue
- **Backend**: Go (Wails)
  - `connector`: Pluggable adapters for different AI agents
  - `discovery`: Dynamic filesystem scanners and configuration parsers
  - `mcpclient`: Minimal MCP-over-stdio client for live server introspection
  - `service`: Business logic and Wails bindings
  - `repository`: Configuration persistence and storage

## Roadmap

- [x] Direct config-file reading and Skill discovery
- [x] Native MCP-over-stdio client for live tool introspection
- [x] Dynamic Agent Discovery and generic connectors
- [ ] Agent detail configuration editor and Settings UI
- [ ] SQLite-backed Profile model for context-aware server grouping
- [ ] Configuration Exporter (sync OwnDeck state back to agent configs)

## Contributing

We welcome contributions from the community. If you'd like to improve OwnDeck, please feel free to submit a Pull Request or open an Issue to discuss proposed changes.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
