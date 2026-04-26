# 多余滚动条问题排查指南

## 核心原则

> **永远不要使用 `overflow-auto`，除非你明确需要两个方向都能滚动。**
> 
> 99% 的场景应该用 `overflow-x-hidden overflow-y-auto`（只允许纵向滚动）或 `overflow-hidden`（完全禁止）。

---

## 常见原因清单

### 1. `overflow-auto` 触发双向滚动条

**问题：** Tailwind 的 `overflow-auto` = CSS `overflow: auto`，同时启用水平和垂直滚动。如果内容宽度比容器多出哪怕 1px，就会出现水平滚动条。

**修复：** 
```css
/* ❌ 错误 */
overflow-auto

/* ✅ 正确 — 只允许纵向滚动 */
overflow-x-hidden overflow-y-auto
```

**已踩坑的组件：**
- `shadcn-vue` 的 `SidebarContent.vue`（默认生成的代码用了 `overflow-auto`）

---

### 2. 嵌套布局的 overflow 冲突

**问题：** 多层容器（html → body → SidebarProvider → SidebarInset → main）都有默认 overflow 行为。如果外层没有 `overflow-hidden`，内层的溢出会"冒泡"到外层显示滚动条。

**修复：** 在 CSS base layer 锁定根元素：
```css
@layer base {
  html {
    @apply overflow-hidden h-full;
  }
  body {
    @apply overflow-hidden h-full;
  }
}
```

让唯一需要滚动的容器（通常是 main content 区域）拥有 `overflow-y-auto`。

---

### 3. 固定宽度容器中的内容撑破

**问题：** sidebar 宽度固定（如 `--sidebar-width: 16rem`），但内部元素（文字、badge、按钮）没有 `truncate` 或 `min-w-0`，导致内容超出触发水平滚动。

**修复：**
```html
<!-- ❌ 文字可能撑破容器 -->
<span>Very Long Navigation Item Name</span>

<!-- ✅ 加 truncate -->
<span class="truncate">Very Long Navigation Item Name</span>
```

对 flex 子元素，加 `min-w-0` 允许 truncate 生效：
```html
<div class="flex items-center min-w-0">
  <span class="truncate">...</span>
</div>
```

---

### 4. padding/margin 导致的 1px 溢出

**问题：** `p-6` (24px padding) 在紧凑布局中可能导致 `content + padding > container`，触发滚动条。

**修复：** 使用 `box-border`（Tailwind 默认）+ 确保容器用 `w-full` 而不是固定宽度。

---

### 5. SidebarProvider `min-h-svh` 与 Wails 窗口

**问题：** `SidebarProvider` 默认使用 `min-h-svh`（100svh），在 Wails 桌面应用中如果内容超出窗口高度，body 会出现纵向滚动条。

**修复：** 确保 `html` 和 `body` 都是 `h-full overflow-hidden`，并且 content 区域用 `flex-1` + `overflow-y-auto` 代替 `min-h-svh`。

---

## 排查流程

```
出现多余滚动条
  ├─ 是水平滚动条？
  │   ├─ 检查容器是否用了 overflow-auto → 改为 overflow-x-hidden overflow-y-auto
  │   ├─ 检查内容是否缺少 truncate / min-w-0
  │   └─ 检查是否有 padding 导致内容超出
  │
  └─ 是纵向滚动条？
      ├─ 检查 html/body 是否缺少 overflow-hidden h-full
      ├─ 检查嵌套容器是否有多个 overflow-auto
      └─ 检查是否有 min-h-svh 导致内容超出视口
```

---

## shadcn-vue 组件已知问题

| 组件 | 问题 | 修复方法 |
|------|------|----------|
| `SidebarContent.vue` | 默认 `overflow-auto` 导致水平滚动条 | 改为 `overflow-x-hidden overflow-y-auto` |
| `SidebarInset.vue` | 缺少 `overflow-hidden`，内容溢出到 body | 添加 `overflow-hidden` |
| `ScrollArea` | 嵌套使用时可能产生双滚动条 | 避免在已有 overflow 的容器内嵌套 ScrollArea |
| `Sheet` / `Dialog` | 打开时可能让 body 出现滚动条 | 组件已处理 body scroll lock，一般不需额外处理 |

---

## 标准页面布局模式

每个页面都应该遵循以下结构，确保只有 content 区域纵向滚动：

```html
<div>  <!-- 页面根 -->
  <PageHeader />  <!-- h-14 shrink-0 固定高度 -->
  <div class="flex-1 overflow-y-auto p-6">  <!-- 可滚动内容区 -->
    <!-- 页面内容 -->
  </div>
</div>
```

**关键点：**
- `PageHeader` 用 `shrink-0` 固定不动
- 内容区用 `flex-1 overflow-y-auto` 填满剩余空间并纵向滚动
- 不要用 `overflow-auto`（会启用双向滚动）
- 不要用 `ScrollArea` 包裹整个内容区域（嵌套滚动容器）
- 不要用 `h-[calc(100vh-xxx)]`（脆弱的计算，flex-1 更可靠）

