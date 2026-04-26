# OwnDeck - Local AI Agent Control Plane 计划书 v0.2

注:OwnDeck 是当前仓库名。产品名可以后续再定,但先不要为了命名停太久。

## 一、是什么

OwnDeck 是一个跑在用户电脑上的 local-first 桌面应用,用于统一管理本机 AI 客户端的工具、Agent、配置和运行环境。

短期定位:

> 给个人开发者的 local-first AI 控制平面: 真正讲 MCP 协议、把 Skills 当一等资产、为 A2A 预留位置——不绑定团队、云和 Kubernetes。

长期定位:

> 本地 AI Agent 控制平面,让本机所有 AI agents 能被发现、被配置、被调用,并逐步通过 A2A 协同工作。

现状判断分两类竞品:

**第一类: MCP 配置同步工具** (mcpm.sh / MCP Linker / Vibe Manager / MCPX / MyMCP)

这是最拥挤的一档,五个项目都在做"一处配置多处同步"。它们的共同短板:

- 都是 JSON / TOML 操作工具,不真正讲 MCP 协议,所以不知道每个 server 暴露了什么工具、参数 schema 是什么、能不能跑起来。
- 只管 server 级别,不管 server 内部的工具粒度。
- 不把 Skills(SKILL.md)当作一等资产。
- 完全不碰 A2A。

**第二类: 全栈控制平面** (Stacklok ToolHive)

ToolHive 已经做 MCP + Skills + A2A 三件事,产品形态最接近 OwnDeck。但它是**企业向**: CLI + 桌面 + 云 + Kubernetes Operator,目标是团队治理。个人用户首页就被 K8s Operator 劝退。OwnDeck 跟它**不直接竞争**,只服务个人 / local-first / 纯桌面这一段——同样三件事,两套人群。

OwnDeck 把上述四件事(MCP 协议探查、Skills 一等公民、A2A 预留、工具级粒度)一次做全,但只服务个人。它的核心是一个统一配置层 + 一个本地探查器:

```text
OwnDeck Profile
  -> Claude Desktop config
  -> Claude Code local/project/user scope
  -> Codex config.toml
  -> Cursor / Zed / VS Code configs
  -> A2A agent registry
```

## 二、为什么做

### 用户痛点

- 用户同时在用 Claude Desktop、Claude Code、Codex、Cursor、Zed、VS Code、自建 AI 工具,每个都有自己的配置格式和配置位置。
- MCP servers 经常要在多个客户端重复配置,JSON / TOML / YAML 格式不同,scope 也不同。
- 配置文件散落在全局目录、项目目录、隐藏目录里,编辑容易出错,也不容易回滚。
- 想知道“这个项目里 Claude 和 Codex 到底能用哪些工具”并不直观。
- A2A 方向正在出现,但个人用户缺少一个本地 agent registry 和调试入口。

### 市场判断

纯 MCP 配置管理已经有人做,包括一键安装、菜单栏启停、多客户端同步等方向。因此 OwnDeck 不应定位为“又一个 MCP GUI”。

仍然有机会的空位是:

- 以项目为中心的 AI runtime profile。
- 跨 Claude / Codex / Cursor / Zed 的配置编译层。
- MCP 工具层与 A2A agent 层的统一资产视图。
- 本地优先、面向个人开发者的 agent control plane。

## 三、产品核心概念

### Profile

Profile 是 OwnDeck 的核心对象,表示某个用户或项目的一套 AI 运行环境。

示例:

- `global-dev`
- `own-deck-project`
- `writing`
- `research`

一个 Profile 中包含:

- 启用哪些 MCP servers
- 哪些 server 是全局级,哪些是项目级
- 哪些客户端需要导出配置
- env 和 secret 如何注入
- 未来可用的 A2A agents
- 未来可用的 Skills / Prompts

### Connector

Connector 表示一个目标客户端的适配器。

OwnDeck 不能完全通用地扫描出“所有 AI 客户端”,因为本地 AI 工具没有统一注册表,也没有统一的 MCP 配置路径、CLI 命令、scope 规则或项目配置文件标准。OwnDeck 应该维护一个 connector registry,每个 connector 自己负责 probe / import / export / validate。

v0.1 优先:

- Claude Desktop
- Claude Code
- Codex
- Cursor

后续:

- Zed
- VS Code / Copilot
- Windsurf
- Cline / Roo Code
- 自定义客户端

### Asset

Asset 是可被配置或调用的 AI 资产。

v0.1:

- MCP server
- MCP tool metadata
- client config target

v0.2+:

- A2A agent
- Agent Card
- Skill
- Prompt
- Workflow / handoff rule

### 两层架构

OwnDeck 内部把功能切成两层,互不耦合:

#### Layer A — OwnDeck 自己的事实来源

- Asset registry: MCP server / Skill / 未来的 A2A agent 统一资产表
- Profile: 一组 asset 的启用集合
- 用户在 OwnDeck UI 里增删改查
- 不依赖任何外部客户端存在,Codex 没装也能建 Profile

#### Layer B — 推给外部客户端

- 把 Profile 编译成 Codex config.toml / Claude .mcp.json / Cursor mcp.json 等
- 写入前 diff 预览、备份、一键回滚
- 只在用户主动"导出"时才碰外部文件

两层职责清晰: Layer A 是 OwnDeck 自己的库,Layer B 是单向的写入管道。

### MCP-as-Client Introspection

OwnDeck 自己实现一个最小的 MCP client。当 Asset 被录入或刷新时,OwnDeck 临时把对应 server 拉起来,完成 initialize 握手并调用 tools/list,把结构化结果存进 Asset 表。

这是 OwnDeck 区别于其他 MCP GUI 的根本动作: 别人只搬运 JSON,OwnDeck 真正知道每个 server 暴露了什么工具、参数 schema、版本、能不能跑起来。

这一步同时解锁:

- 工具级开关 (在 Profile 里禁用 github 的 delete_repo,而不是禁用整个 server)
- Server 健康检查 (initialize 失败 = 配置坏了)
- 跨客户端版本一致性检测 (同名 server 在 Codex 和 Cursor 是不是同一份)
- 写入外部配置前的预检 (推过去之前先确认它能跑)

OwnDeck 不会常驻持有这些连接,也不会调 tools/call,它只是临时探查、存档、断开。

### MCP 和 Skills 是上下两层,不是替代关系

```text
Skills 层(知识 / 流程)         "怎么用 github 工具做发版"
       ↓ reference
MCP 层(能力 / 工具)            github_create_pr / playwright_click / ...
```

2026 年生态明显往 Skills 倾斜——Anthropic 把 Skills 推成一等公民,Vercel 出 Skills.sh,Codex 插件市场里大部分新插件都是 Skills 而不是 MCP server。原因很直接: 写 Skills 比写 MCP server 成本低 10 倍,通用 MCP server 又已经被大公司做得差不多了。

但 **Skills 不是替代 MCP**,Skills 引用的工具仍然来自 MCP server。OwnDeck 必须**同时管这两层**才完整: Skills 是用户感知到的入口,MCP 是 Skills 跑起来的前提。

竞品在这件事上的位置:

- MCP 同步工具 (Linker / mcpm 等): 只管 MCP 层, 看不到 Skills,2026 年的生态它们看不到。
- Skills 工具 (Skills.sh / Skilldex): 只管 Skills 层, 不管 MCP, Skills 引用的工具它们装不了。
- OwnDeck: 同时管两层,这是真正的差异化。

### 不做 marketplace,做 brew / npm

> OwnDeck 的卖点不是"我这里有 600 个 server",是"你从哪听说的 server,OwnDeck 帮你正确装到所有客户端"。

类比:

- HomeBrew 没有自建商店,它是 install 工具——你说装啥它去对应源拉。
- npm 不教你"该装什么",它是 install 工具——你说装啥它给你装。

OwnDeck 该是 AI 工具的 brew/npm,不是 App Store。原因:

- 自建 marketplace 是 MCP Linker 的赌注 (它砸了 1 年做到 600+ server),OwnDeck 在 marketplace 这条赛道**赢不了**。
- 跟 Codex 自己的插件市场正面竞争更没意义: Codex plugin 是 OpenAI 的封闭格式,装到 Cursor / Claude 上跑不起来。
- 真正没人做对的是"装到对的地方 + 装对了"——写进 Codex config.toml + Claude .mcp.json + Cursor mcp.json,写之前 introspection 一次确认能跑。

未来的搜索/发现入口直接接[官方 MCP Registry](https://registry.modelcontextprotocol.io/),不重复造目录,只做"在 OwnDeck 内一键安装"的 UX。

### 不通过 CLI 抓取

OwnDeck 不调 codex mcp list / claude mcp list 之类的子命令解析输出。每个 connector 直读对应 client 的配置文件 (TOML / JSON),写入也走同一条路径。CLI 输出格式不稳定、字段不全、依赖 CLI 安装,不能作为事实来源。

## 四、v0.1 核心功能

v0.1 的目标不是做完整 A2A 平台,而是先做出“统一配置层”。这是刚需入口,也是未来 A2A 控制平面的地基。

| 功能 | 说明 | 优先级 |
| --- | --- | --- |
| Profile 管理 | 创建、编辑、复制、删除 AI runtime profile | P0 |
| MCP 资产清单 | 统一列出已发现和手动添加的 MCP servers | P0 |
| 直读式自动发现 | 直接解析 Claude Desktop / Claude Code / Codex / Cursor 的配置文件 (TOML / JSON),不通过 CLI 抓取 | P0 |
| MCP 协议探查 | OwnDeck 作为 MCP client 临时连接 server,拉取 tools / prompts / resources 列表存档 | P0 |
| 配置编译 | 将 OwnDeck Profile 导出为不同客户端需要的 JSON / TOML 配置 | P0 |
| 配置备份 | 修改客户端配置前自动备份原文件 | P0 |
| 回滚 | 出错时恢复到上一个可用配置 | P0 |
| 手动添加 | 添加 stdio / HTTP / SSE 类型 MCP server | P0 |
| Skill 资产管理 | 录入、查看、跨客户端复用 SKILL.md 资产 | P0 |
| 健康检查 | 检查 server 是否可启动、是否能握手、是否能列出 tools | P1 |
| 工具级开关 | 在 Profile 中按工具粒度启用 / 禁用 (基于 introspection 拿到的 tools/list) | P1 |
| 启用/禁用 | 在 Profile 中启用或禁用某个 server,不删除原始记录 | P1 |
| 详情查看 | 查看 command、args、env、tools、来源客户端 | P1 |
| 跨客户端一致性 | 同名 server 在不同客户端的版本 / 配置差异提示 | P1 |
| 标签/搜索 | 按用途、项目、风险级别搜索资产 | P2 |

## 五、v0.1 不做什么

- 不做完整 Agent 编排 / Workflow。
- 不做云同步。
- 不做团队协作。
- 不做企业级权限、审计。
- 不做 marketplace。
- 不承诺支持所有 MCP 客户端。
- 不把 A2A 做成主流程,但数据模型要为 A2A 留位置。
- 不直接存储明文 secret,优先使用系统 Keychain / Credential Manager。

## 六、技术栈

- 前端:Vue 3 + TypeScript + shadcn-vue + Tailwind CSS
- 后端:Go + Wails
- 存储:SQLite,优先使用纯 Go 驱动
- 配置解析:JSON / TOML / YAML 结构化解析(BurntSushi/toml 等),不做字符串拼接,不通过 CLI 抓取
- MCP 客户端: 内置最小 MCP-over-stdio client(initialize / tools/list / prompts/list / resources/list),优先复用官方 Go SDK
- Secret:macOS Keychain / Windows Credential Manager / Linux Secret Service
- 协议:v0.1 以 MCP 消费、探查和导出为主,v0.2 开始加入 A2A agent registry
- 打包:Wails,优先 macOS,再 Windows,最后 Linux

预计应用大小:15-40 MB

预计内存占用:50-120 MB

## 七、数据模型草案

核心表:

- `profiles`:运行环境配置集合
- `assets`:MCP server、A2A agent、Skill、Prompt 等统一资产
- `profile_assets`:Profile 与资产的启用关系
- `client_targets`:Claude Desktop、Claude Code、Codex、Cursor 等目标客户端
- `config_sources`:从哪里发现的配置,以及原始路径
- `export_history`:每次导出、备份、回滚记录
- `health_checks`:健康检查结果

关键设计原则:

- MCP server 和 A2A agent 都是 asset,但 type 不同。
- Profile 不直接等于某个客户端配置,而是 OwnDeck 的中间表示。
- 各客户端 connector 只负责 import / export / validate。
- 所有写入外部配置文件的动作必须可备份、可预览、可回滚。

## 八、用户画像

主要用户:

- 同时使用 Claude Code、Codex、Cursor、Zed 的个人开发者
- 重度 MCP 用户
- 自托管和 local-first 爱好者
- 正在尝试 A2A / multi-agent 的开发者

典型场景:

> 我在一个项目里希望 Claude Code、Codex 和 Cursor 都能用同一组 GitHub、Playwright、filesystem、database MCP。现在我要分别改 JSON、TOML 和项目配置。OwnDeck 让我建一个 Profile,勾选目标客户端,然后自动生成、备份、导出和验证配置。

长期场景:

> 我有几个本地 agents:代码审查 agent、测试 agent、文档 agent。OwnDeck 能发现它们的 Agent Card,知道它们会什么,并让 Claude / Codex 通过 MCP bridge 或 A2A 调用它们。

## 九、和现有方案的差异

| 方案 | 核心形态 | 强项 | OwnDeck 差异 |
| --- | --- | --- | --- |
| **Stacklok ToolHive** | CLI + 桌面 + 云 + K8s Operator | 同时做 MCP + Skills + A2A,产品形态最像 OwnDeck | 企业向(团队 / K8s / 治理),个人用户进不去。OwnDeck 是同样三件事的**个人 / local-first 版**,不直接竞争 |
| MCP Linker | Tauri 桌面 app + 600+ server marketplace | 多客户端 JSON 同步、安装一键化 | 不讲 MCP 协议、不管 Skills、无 A2A、无工具级控制 |
| mcpm.sh | CLI + Profile + smart router | Profile 概念成熟、router 共享会话 | CLI 形态、不管 Skills 和 A2A、不做协议级探查 |
| Vibe Manager | 桌面 app | 一处编辑 sync 多端 | 纯 JSON 同步器,没有 introspection / Skills / A2A |
| MCPX | CLI | 多客户端同步 | 跟 mcpm 类似,无 introspection / Skills / A2A |
| MyMCP | macOS app | server 注册表浏览、按 client 加密存 secret | 安装器形态,不管 Skills 和 A2A |
| Skills.sh (Vercel) | CLI + 中央 registry | “npm for AI skills”,标准化 SKILL.md 分发 | 只管 Skills,不管 MCP 基础设施 |
| Skilldex (skillpm) | CLI + Hono/Supabase registry | 三层 scope、human-in-the-loop | 学术原型,只管 Skills |
| MCP Gateway Registry | 企业级 gateway | OAuth / 治理 / A2A | 企业部署,不是本地优先,个人开发者用不上 |
| MCP Inspector (官方) | 调试工具 | 实时连接 server 做协议调试 | 一次性调试器,不存档、不管多 server、不做 Profile |
| getmcp / Smithery / 官方 Registry | 网页目录 | 发现和安装 server | 网页目录,不管本机运行环境 |

核心差异:

> ToolHive 之外其他工具都是”单层搬运工”——要么只搬 MCP JSON,要么只发 Skills 包。ToolHive 是全栈控制平面但企业向。OwnDeck = ToolHive 的个人 / local-first 版本: 同时管 MCP / Skills / A2A 三层,但只服务单个开发者本机,不绑定团队和云。一句话: **它们做的是单层同步或企业治理,OwnDeck 做的是个人本机控制平面**。

## 十、开发计划

### 阶段 0:命名和地基

时长:2-3 天

产出:

- 确认 OwnDeck 作为临时产品名
- 初始化 Wails + Vue + TypeScript 项目
- 建立 README、PROGRESS.md、基本开发命令

### 阶段 1:直读式配置发现 ✅ 基本完成 (2026-04-26)

时长:3-5 天

产出:

- ✅ 删除 codex mcp list / claude mcp list 子进程抓取逻辑
- ✅ 用 BurntSushi/toml 直读 Codex `~/.codex/config.toml`
- ✅ **额外**: 遍历 Codex 启用的插件目录,读 `.codex-plugin/plugin.json` → `.mcp.json`,把插件捎带的 MCP server 也收进来。这件事 Codex 自己的桌面 UI 都没做对(参见 [openai/codex#17360](https://github.com/openai/codex/issues/17360))
- ✅ 给每个 server 加 `Origin` 字段 (`user` / `plugin:NAME@MARKETPLACE`),UI 可分组显示
- ✅ 直读 Claude Code 的 `~/.claude.json`、`~/.claude/settings.json`、项目 `.mcp.json`
- ✅ 直读 Claude Desktop 的 `~/Library/Application Support/Claude/claude_desktop_config.json`
- ✅ Antigravity 早就直读 JSON,不变
- ⏳ Cursor connector (格式标准 JSON,与 Antigravity 同源,10 分钟工作量,等用户装了再补)

验收:

- ✅ 不写任何外部配置文件
- ✅ CLI 没装也能展示配置
- ✅ 在自己的机器上准确展示了 1 个 plugin-bundled MCP server + 13 个 Skills

工程副产物:

- 后端分层重构: `cmd / internal/{app,connector,service,repository,platform,discovery}` + Connector 接口 + Registry,DI 在 main.go 一处装配
- 前端踩坑文档: `docs/frontend-pitfalls.md` (v-for / v-if 兄弟节点白屏、Go nil slice 跨 Wails 变 JSON null 两条)

### 阶段 2:MCP 协议探查

时长:1-2 周

产出:

- internal/mcpclient: 最小 MCP-over-stdio client (优先复用官方 Go SDK)
- 实现 initialize / notifications/initialized / tools/list / prompts/list / resources/list
- Asset enrichment: server 录入或刷新时,临时连接、拉清单、存档、断开
- 健康状态字段 (启动失败 / 协议失败 / 正常)
- 进程生命周期管理 (超时、并发上限、错误恢复)

验收:

- 对常见 stdio server (filesystem / github / playwright) 能拿到结构化工具清单
- initialize 失败的 server 在 UI 上明确标红
- 不会因为某个 server 卡住影响整体探查

### 阶段 3:OwnDeck Profile + Asset

时长:1-2 周

产出:

- SQLite schema (profiles / assets / profile_assets / config_sources / health_checks)
- Profile CRUD
- MCP / Skill asset CRUD
- 工具级 enable / disable (基于阶段 2 的 tools/list)
- 导入现有发现结果为 Profile

验收:

- 能从已发现的配置生成一个 OwnDeck Profile,并按工具粒度调整
- Skill 和 MCP server 在同一个 Profile 视图里共存

### 阶段 4:配置编译和导出

时长:2-3 周

产出:

- Codex TOML exporter (优先,格式相对简单)
- Claude Code / Claude Desktop exporter
- Cursor exporter
- 写入前 preview diff
- 自动备份 + 一键回滚
- 推送前用阶段 2 的 introspection 做预检 (server 拉得起再写入)

验收:

- 建一个 Profile,能导出到至少 3 个客户端
- 出错能回滚
- 推坏配置之前能在 diff 视图里看到风险

### 阶段 5:打磨发布

时长:2 周

产出:

- 跨客户端版本一致性提示
- 错误提示和修复建议
- README、demo GIF、首发 blog

验收:

- 完整跑通”直读 -> 探查 -> 建 Profile -> 工具级编辑 -> 导出 -> 回滚”
- macOS 稳定运行
- Windows 至少跑通核心流程

总计:

- 全职:8-12 周
- 业余:3-4 个月

## 十一、A2A 路线图

v0.1:Unified Profile Compiler

- 跨客户端 MCP 配置统一层
- 为 asset model 预留 `a2a_agent` 类型

v0.2:Agent Registry

- 手动添加 A2A endpoint
- 拉取和展示 Agent Card
- 健康检查
- 本地 agent 目录

v0.3:MCP <-> A2A Bridge

- 将 A2A agents 暴露成 MCP tools
- 让 Claude / Codex 等 MCP 客户端可以调用 A2A agent
- 支持基础 message send、artifact 查看

v0.5:Local Agent Control Plane

- Agent capability search
- 简单 handoff rule
- 本地任务路由
- 项目级 agent graph

v1.0:

- 可选云同步
- 团队 Profile
- 权限和审计
- 更完整的 agent workflow

## 十二、发布策略

首发不要讲太大的”AI agent 操作系统”。已经有一堆工具在讲”一处配置多处同步”,再讲一遍没意义。新的一句话定位:

> The only local-first console that actually speaks MCP — manage tools, skills, and (soon) A2A agents in one place.

中文版:

> 不只是同步 JSON。OwnDeck 真正讲 MCP 协议,看得到每个 server 暴露的工具,管得了 Skills,未来还有 A2A。

### 首发 demo 的具体钩子

抽象的”全栈控制平面”传播力弱,首发要选一个 5 秒就能炸出来的视觉 wedge。候选(按可演示性排序):

| Wedge | 5 秒 demo | 现有竞品 |
| --- | --- | --- |
| **看见每个 MCP server 内部的工具** | 一张截图: 左边”github 服务器”,右边展开 23 个工具,每个有 description / 参数 schema | 几乎没人做 GUI |
| **跨客户端 Skills 总览** | “13 个 Skills 来自 Codex / Antigravity”,MCP Linker 等连 Skills 概念都没有 | 没人做 |
| **Codex 插件 server 显示对了** | “OwnDeck 显示 6 个,Codex 自己的 GUI 显示 0 个”——直接对着 [openai/codex#17360](https://github.com/openai/codex/issues/17360) 链过去 | Codex 自己都没修 |
| **工具级开关** | 禁用 github 的 delete_repo,其他工具不动 | 需要 introspection 先做 |

首发 blog 优先用前两条,因为 “introspection 看见工具” 和 “Skills 视图” 都已经在阶段 1-2 内,做完直接出片。第三条作为切入点段落: “为什么 OwnDeck: 因为 Codex 自己都没把这件事做对”。

渠道:

- Hacker News:Show HN
- Reddit:r/LocalLLaMA、r/selfhosted、r/ClaudeAI、r/Cursor
- GitHub README
- Twitter/X AI developer 圈
- V2EX、即刻、少数派

发布材料:

- 30 秒 demo:一个 Profile 同步到 Claude Code + Codex + Cursor
- README GIF
- 配置格式对比图
- Blog:《为什么 AI 客户端需要一个统一配置层》

## 十三、风险

| 风险 | 应对 |
| --- | --- |
| MCP 配置同步类工具已经红海 (Linker / mcpm / Vibe Manager / MCPX / MyMCP) | 不跟它们正面拼”同步”,而是把 MCP 协议探查、Skills、A2A、工具级粒度做全,用差异化压制 |
| 客户端配置格式变化 | 每个客户端独立 connector,加版本检测和备份回滚 |
| MCP server 启动卡死 / 副作用 | introspection 有超时上限、并发上限,临时连接立刻断开,不持有长连接 |
| 写坏用户配置 | 默认只读发现,写入前 diff,自动备份,一键回滚,推送前用 introspection 预检 |
| A2A 生态太早 | v0.1 不依赖 A2A 成熟,只预留模型和路线;v0.2 才接入 |
| 跨平台打包麻烦 | 先 macOS,再 Windows,Linux 后置 |
| Secret 管理风险 | 不存明文 secret,使用系统凭据管理 |
| 用户看不懂”控制平面”愿景 | 首发的 demo 直接展示”看到 server 内部工具 + 工具级开关”,这是肉眼可见的差异化 |

## 十四、成功标准

### v0.1 发布时

- 自己每天用 OwnDeck 管理 Claude / Codex / Cursor 配置。
- 能完整跑通“发现 -> Profile -> 导出 -> 健康检查 -> 回滚”。
- 至少支持 Claude Desktop、Claude Code、Codex、Cursor 中的 3 个。
- 修改外部配置前 100% 自动备份。
- README 能让陌生开发者 5 分钟内理解项目价值。

### v0.1 后 3 个月

- 20+ 真实安装用户。
- 5+ 日常活跃用户。
- 10+ 来自真实使用的 issue。
- 1+ 外部 contributor。
- 明确 v0.2 A2A Agent Registry 是否值得继续。

## 十五、现在就开始

第一周不要做华丽 UI。当前目标只有一个:

> 把所有 client 的发现路径从 CLI 抓取换成直读配置文件,然后立刻让 OwnDeck 自己讲一次 MCP 协议。

最小任务:

1. 删除现有 codex / claude CLI 抓取代码,connector 改成直接解析 TOML / JSON。
2. 落地最小 MCP-over-stdio client (initialize + tools/list)。
3. 让发现出的每个 server 都尝试 introspect 一次,UI 展示工具清单和健康状态。
4. UI 展示:客户端、配置路径、server 名、类型、command/url、env、tools 清单、健康。
5. 写 `PROGRESS.md`,每周记录一次进展。

不要一开始做 marketplace、A2A workflow、云同步、复杂视觉设计。先让 OwnDeck 成为你自己电脑上的真实工具,而且是唯一一个真正"看得见"每个 server 的工具。
