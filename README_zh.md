<div align="center">
  <img src="build/appicon.png" width="128" alt="OwnDeck Logo">
  <h1>OwnDeck</h1>
  <p><b>本地优先的 AI Agent 控制中心</b></p>

  <a href="https://wails.io/"><img src="https://img.shields.io/badge/Wails-v2-blue?style=for-the-badge&logo=go" alt="Wails"></a>
  <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/Vue.js-3.0-4FC08D?style=for-the-badge&logo=vue.js" alt="Vue 3"></a>
  <a href="https://www.typescriptlang.org/"><img src="https://img.shields.io/badge/TypeScript-Ready-3178C6?style=for-the-badge&logo=typescript" alt="TypeScript"></a>
  <br><br>

  [English](README.md) | [简体中文](README_zh.md)
</div>

<hr>

OwnDeck 是首个真正原生支持 MCP (Model Context Protocol) 协议的**本地优先 (Local-first)** 控制台。它能在一个统一的界面中管理所有的 MCP Server（精确到具体工具级别）以及各种 AI Agent，是你本地 AI 开发的终极控制中心。

市面上的同类工具往往只做简单的 JSON 配置文件同步，而 OwnDeck 具备**协议级内省 (Introspection)** 能力，能够理解每个 Server 暴露的工具，并将 Skills 和 Agents 视为与 MCP 同等重要的一等资产。

## ✨ 核心特性

- **🌐 动态 Agent 发现**: 自动扫描文件系统，精准探测已安装的 AI Agent（如 Claude, Codex, Gemini CLI, Cursor）及其关联扩展。
- **🔌 协议级内省**: 直接通过 `stdio` 与 MCP Server 通信，实时查询可用的 Tools、Prompts 和 Resources。
- **📖 直读式配置解析**: 告别脆弱的 CLI 子进程抓取，OwnDeck 原生解析 JSON/TOML 配置文件，零外部依赖。
- **🛠️ 集中式控制台**: 在一个现代化的可视化面板中，统览全机的所有 MCP Server、连接健康状态及工具清单。
- **⚡ 无感自动刷新**: 基于 Vue 3 + shadcn-vue + Wails 构建的现代响应式栈，支持静默后台轮询与无缝状态管理。
- **🌍 原生多语言**: 开箱即用的中英文双语支持。

## 🚀 快速开始

### 环境依赖
- Go 1.21+
- Node.js & npm
- [Wails v2](https://wails.io/)

### 安装步骤

1. **安装 Wails CLI**:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   wails doctor
   ```

2. **克隆代码库**:
   ```bash
   git clone https://github.com/loqbit/OwnDeck.git
   cd OwnDeck
   ```

3. **运行开发模式**:
   ```bash
   wails dev
   ```

4. **构建生产包**:
   ```bash
   wails build
   ```

## 🏗️ 架构设计

OwnDeck 采用了面向扩展的分层架构：

- **前端**: Vue 3 + TypeScript + Tailwind CSS + shadcn-vue
- **后端**: Go (Wails)
  - `connector`: 适配各种 AI Agent 的可插拔连接器
  - `discovery`: 动态文件系统扫描器与 MCP 配置解析器
  - `mcpclient`: 轻量级 MCP-over-stdio 客户端，用于进程启动与协议内省
  - `service`: 业务聚合逻辑与 Wails 接口绑定
  - `repository`: 配置持久化与本地存储

## 🗺️ 演进路线图

- [x] **阶段 1**: 直读式配置发现与 Skill 资产识别
- [x] **阶段 2**: 原生 MCP-over-stdio 客户端与实时工具内省
- [x] **阶段 3 (第一期)**: 动态 Agent 扫描与泛型扩展支持
- [ ] **阶段 3 (第二期)**: Agent 详细配置编辑器与设置 UI
- [ ] **阶段 4**: 基于 SQLite 的 Profile 管理模型（按业务场景打包 MCP Server）
- [ ] **阶段 5**: 配置编译器与反向导出器（一键将配置同步至各大客户端）

## 🤝 参与贡献
欢迎随时提交 Pull Request 或 Issue！

## 📄 许可证
本项目采用 MIT 许可证 - 详情请见 [LICENSE](LICENSE) 文件。
