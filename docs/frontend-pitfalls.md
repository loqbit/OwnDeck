# Frontend Pitfalls

OwnDeck 前端踩过的坑和复盘。新坑往最前面加，老坑别删（搜索方便）。

每条记录的格式：

- **症状**：肉眼或控制台看到了什么
- **触发**：什么操作触发的
- **根因**：为什么会出
- **修法**：这次怎么修的
- **预防**：以后写代码该遵守什么规则

---

## 2026-04-26 — Go 的 nil slice 跨 Wails 变成 JSON `null`，前端拿到后炸列表

- **症状**：报错 `TypeError: null is not an object (evaluating '$setup.filteredServers.length')`，调用栈在 `runtime-core.esm-bundler.js`。修了之前那条 `v-for/v-if` 兄弟坑后还是炸。注意这次错的是 `filteredServers.length` 整体——computed 本身返回了 null，不是它里面的某个 item。
- **触发**：连接一个返回零 server 的 client（比如刚装 Claude Desktop，还没配 mcpServers），切到 MCP 服务器页。
- **根因**：Go 这样写：
  ```go
  var results []T   // 这是 nil，不是空切片
  return results, nil
  ```
  `encoding/json` 把 nil slice 序列化成 `null`，**不是 `[]`**。Wails 把这个 `null` 透传到前端，赋值给 `servers.value` 后，`filteredServers` computed 返回 `servers.value` 也是 null，模板里 `.length` / `.filter` 都炸。
  之所以这次才炸：之前 connector 都是 `make([]T, 0, n)` 起手所以非 nil，但 [registry.aggregate](../internal/connector/connector.go) 用的是 `var results []T`，当所有 connector 都返回 0 项时 `results` 永远没被 append 过，保持 nil。
- **修法**：两层都要修。
  - 后端：用 `make([]T, 0)` 起手，绝不返回 nil slice 给前端。改在 [connector.aggregate](../internal/connector/connector.go) 里。
  - 前端：所有 `await` 后端方法的地方都加 `?? []` 兜底（`useMCPServers` / `useSkills` / `useClients`）。后端再保证非 nil，前端这层是为了防新加的方法忘记。
- **预防**：
  - **Go 侧任何返回 slice 给前端的函数，必须 `make([]T, 0)` 起手，不要用 `var x []T`**。后端约定一致就不会出错。
  - **前端侧每个 `await Backend.Foo()` 的赋值都 `?? []` 或 `?? {}`**。Wails 的 TS 类型不会告诉你后端可能返回 null（因为 Go 类型上写的是 `[]T` 而不是 `[]T | nil`，但跨边界 nil 切片就是 null）。
  - **响应式 `ref`、`computed` 永远默认提供初始值**：`ref<T[]>([])` 而不是 `ref<T[] | null>(null)`。
  - **空状态、loading 态、错误态三件套**还是要写。

---

## 2026-04-26 — TableBody 里 `v-for` 和 `v-if` 兄弟节点交替导致白屏

- **症状**：点击「MCP 服务器」标签页时整页变白；DevTools 控制台报 `Unhandled Promise Rejection: TypeError: null is not an object (evaluating 'node.parentNode')`，调用栈在 `runtime-core.esm-bundler.js` 里。
- **触发**：连接 Claude Desktop 之后切换到 MCP 服务器页。connected client 列表变化导致 `filteredServers` 在「有数据」和「空」之间快速切换。
- **根因**：模板里这样写：
  ```vue
  <TableBody>
    <TableRow v-for="s in filteredServers" :key="..."> ... </TableRow>
    <TableEmpty v-if="filteredServers.length === 0" ...> ... </TableEmpty>
  </TableBody>
  ```
  `TableEmpty` 内部本身也是个 `<TableRow>`。当列表从空变成非空（或反过来），Vue 的 patch 算法要同时插入 v-for 出来的兄弟节点 + 移除 v-if 的 `TableEmpty`。某些时机下被移除节点的 `parentNode` 已经是 null，于是炸掉整个组件树，触发上面的报错并让整页白屏。
- **修法**：改成 `v-if / v-else` 互斥：
  ```vue
  <TableBody>
    <template v-if="filteredServers.length > 0">
      <TableRow v-for="s in filteredServers" :key="..."> ... </TableRow>
    </template>
    <TableEmpty v-else ...> ... </TableEmpty>
  </TableBody>
  ```
  顺带把 key 加上 `sourcePath` 段，避免未来同名 server 撞 key。
- **预防**：
  - **同一个父节点里不要同时出现 `v-for` 列表和一个会替代它的 `v-if` 兄弟**——尤其当兄弟也是同类型节点（比如 `TableRow`）。永远用 `v-if / v-else` 互斥分支，每次只渲染一种。
  - **`v-for` 的 `:key` 要全局唯一**。`name` 不够用就加来源路径。Codex 直读后同名 server 可能从用户配置和插件两处都来，必须 `${client}-${name}-${sourcePath}` 这种组合 key。
  - **空状态尽量放到 `<Table>` 外面**，不要塞进 `TableBody`。`TableBody` 只放数据行最稳。

---

## 通用排查清单

下次再白屏 / 渲染异常，按这个顺序查：

1. **DevTools 控制台**（应用窗口右键 → Inspect → Console）
   - 红色报错的第一行才是真线索，下面的调用栈用来定位组件
   - `null is not an object (evaluating 'node.xxx')` → Vue 渲染 patch 出错，去看 v-for / v-if / key
   - `Cannot read properties of undefined` → 数据没到位就被模板访问了，加可选链 `?.` 或默认值
2. **Network 面板**
   - 不是用来看 wails IPC 的（IPC 走 `window.go.*`，不是 HTTP）
   - 主要看 vite 是不是有 module 加载失败
3. **`wails dev` 终端日志**
   - Go 侧 panic 会打在这里，前端只会表现成 IPC 调用挂起
   - bindings 重新生成的日志也在这里，看是不是真的更新了
4. **检查 `frontend/wailsjs/go/` 是不是最新**
   - 改了 Go 类型字段后，必须等 wails dev 重新生成
   - 必要时 `Ctrl+C` 重启 wails dev，热重载偶尔会卡

## 写前端时的几条铁律

- **`v-for` 和 `v-if` 不能在同一个元素上**，Vue3 已经报错；放兄弟节点也要警惕本文档第一条的坑。
- **列表 `:key` 必须真正唯一**。一旦未来数据可能撞名，提前把唯一性维度（路径 / id / origin）加进去。
- **后端 struct 字段加新的，前端可以不动**——TS 类型只是描述，运行时多出来的字段是 noop。
- **后端 struct 字段改名 / 删除**，前端必须同步检查所有 `source["xxx"]` / 模板里 `server.xxx` 的引用。
- **空状态、loading 态、错误态三件套永远要写**，否则数据没到位时模板就会去访问 undefined 的字段。
