<div align="center">
  <img src="build/appicon.png" width="128" alt="OwnDeck Logo">
  <h1>OwnDeck</h1>
  <p><b>The Local-First AI Agent Control Plane</b></p>

  <a href="https://wails.io/"><img src="https://img.shields.io/badge/Wails-v2-blue?style=for-the-badge&logo=go" alt="Wails"></a>
  <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/Vue.js-3.0-4FC08D?style=for-the-badge&logo=vue.js" alt="Vue 3"></a>
  <a href="https://www.typescriptlang.org/"><img src="https://img.shields.io/badge/TypeScript-Ready-3178C6?style=for-the-badge&logo=typescript" alt="TypeScript"></a>
  <br><br>

  [English](README.md) | [简体中文](README_zh.md)
</div>

<hr>

OwnDeck is the **only local-first console** that actually speaks the MCP (Model Context Protocol). It manages MCP servers (down to the individual tool) and AI agents in one place, serving as your ultimate command center for local AI development.

Unlike existing tools that merely synchronize JSON config files, OwnDeck **introspects** each server, understands the tools it exposes, and treats Skills and Agents as first-class assets.

## ✨ Features

- **🌐 Dynamic Agent Discovery**: Automatically scans your filesystem to detect installed AI agents (Claude, Codex, Gemini CLI, Cursor) and their extensions.
- **🔌 Protocol-Level Introspection**: Directly communicates with MCP-over-stdio servers to live-query tools, prompts, and resources. 
- **📖 Direct Config Reading**: No fragile CLI subprocess scraping. OwnDeck reads configurations natively (JSON/TOML) with zero external dependencies.
- **🛠️ Centralized Control**: View all your MCP servers, their connection health, and available tools across your entire machine in one beautiful dashboard.
- **⚡ Auto-Refresh UI**: Built with a reactive, modern stack (Vue 3 + shadcn-vue + Wails) featuring silent background polling and smooth state management.
- **🌍 Internationalization**: First-class English and Simplified Chinese support out of the box.

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Node.js & npm
- [Wails v2](https://wails.io/)

### Installation

1. **Install Wails CLI**:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   wails doctor
   ```

2. **Clone the repository**:
   ```bash
   git clone https://github.com/loqbit/OwnDeck.git
   cd OwnDeck
   ```

3. **Run in Development Mode**:
   ```bash
   wails dev
   ```

4. **Build for Production**:
   ```bash
   wails build
   ```

## 🏗️ Architecture

OwnDeck uses a layered architecture designed for extensibility:

- **Frontend**: Vue 3 + TypeScript + Tailwind CSS + shadcn-vue
- **Backend**: Go (Wails)
  - `connector`: Pluggable adapters for different AI agents (Claude, Codex, etc.)
  - `discovery`: Dynamic filesystem scanners and MCP config parsers
  - `mcpclient`: Minimal MCP-over-stdio client for live server introspection
  - `service`: Business logic and Wails bindings
  - `repository`: Configuration persistence and storage

## 🗺️ Roadmap

- [x] **Stage 1**: Direct config-file reading and Skill discovery
- [x] **Stage 2**: Native MCP-over-stdio client for live tool introspection
- [x] **Stage 3 (Phase 1)**: Dynamic Agent Discovery and generic connectors
- [ ] **Stage 3 (Phase 2)**: Agent detail configuration editor and Settings UI
- [ ] **Stage 4**: SQLite-backed Profile model (group MCP servers into functional profiles)
- [ ] **Stage 5**: Configuration Compilation + Exporter (write back to client configs)

## 🤝 Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
