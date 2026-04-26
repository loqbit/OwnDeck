<h1 align="center">OwnDeck</h1>

<p align="center">
  <strong>本地优先的 AI Agent 控制中心</strong>
</p>

<p align="center">
  <a href="https://wails.io/"><img src="https://img.shields.io/badge/Wails-v2-blue?style=flat-square&logo=go" alt="Wails"></a>
  <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/Vue.js-3.0-4FC08D?style=flat-square&logo=vuedotjs" alt="Vue 3"></a>
  <a href="https://www.typescriptlang.org/"><img src="https://img.shields.io/badge/TypeScript-Ready-3178C6?style=flat-square&logo=typescript" alt="TypeScript"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square" alt="License"></a>
</p>

<p align="center">
  <a href="#核心特性">核心特性</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#架构设计">架构设计</a> •
  <a href="#参与贡献">参与贡献</a>
</p>

<p align="center">
  <em><a href="README.md">English</a> • <a href="README_zh.md">简体中文</a></em>
</p>

---

OwnDeck 是专为 MCP (Model Context Protocol) 打造的本地控制台。它将你的 MCP Servers、Tools、Skills 与各大 AI Agent 聚合在统一的可视化面板中，构筑你本地 AI 开发的终极指挥中心。

市面上的同类工具大多只做简单的 JSON 同步，而 OwnDeck 具备**协议级内省 (Introspection)** 能力，通过与服务器建立 `stdio` 连接，精准探测其实际暴露的工具与能力边界。

## 核心特性

- **动态 Agent 发现**: 自动扫描文件系统，精准识别并纳管本机安装的 AI Agent（Claude, Codex, Gemini CLI, Cursor）及其扩展插件。
- **协议级内省**: 原生发起 MCP 协议握手，实时探测 Server 侧的 Tools、Prompts 和 Resources，告别黑盒。
- **直读式配置解析**: 零外部依赖，直接解析底层 JSON/TOML 配置，彻底摒弃脆弱的 CLI 子进程数据抓取方案。
- **全景仪表盘**: 现代化的 UI 界面，一览所有 MCP Server 的分布、连接健康度以及细粒度工具清单。
- **高性能前端架构**: 基于 Vue 3 + Wails 构建，底层利用静默后台轮询机制，实现无感的 UI 状态热更新。

## 快速开始

### 环境准备

确保你的系统已安装以下环境：
- [Go](https://golang.org/doc/install) 1.21 或更高版本
- [Node.js](https://nodejs.org/) & npm
- [Wails v2](https://wails.io/docs/gettingstarted/installation)

### 安装指引

1. 安装 Wails CLI 并校验环境：
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   wails doctor
   ```

2. 克隆项目仓库：
   ```bash
   git clone https://github.com/loqbit/OwnDeck.git
   cd OwnDeck
   ```

3. 启动开发模式：
   ```bash
   wails dev
   ```

4. 构建生产应用：
   ```bash
   wails build
   ```

## 架构设计

OwnDeck 采用逻辑解耦的清晰分层架构，确保协议发现层与视图层完全分离：

- **前端层**: Vue 3 + TypeScript + Tailwind CSS + shadcn-vue
- **后端层**: Go (Wails)
  - `connector`: 适配异构 AI Agent 的可插拔连接器
  - `discovery`: 动态文件系统探针与 MCP 配置解析引擎
  - `mcpclient`: 极简 MCP-over-stdio 客户端，负责唤起进程与协议握手
  - `service`: 业务逻辑聚合与状态组装
  - `repository`: 核心配置的本地持久化与序列化

## 演进路线

- [x] 直读式配置发现与 Skill 资产识别
- [x] 原生 MCP-over-stdio 客户端与实时工具内省
- [x] 动态 Agent 扫描与泛型连接器支持
- [ ] Agent 详细配置编辑器与高级设置面板
- [ ] 基于 SQLite 的 Profile 管理模型（按业务场景隔离 MCP Server）
- [ ] 核心配置反向导出器（一键将沙盒配置同步回各大客户端）

## 参与贡献

我们非常欢迎来自社区的贡献。如果您有好的想法或发现了问题，请随时提交 Pull Request 或发起 Issue 进行讨论。

## 许可证

本项目采用 Apache License 2.0 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。
