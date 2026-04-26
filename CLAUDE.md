# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.

---

## 项目特定规则 (OwnDeck)

### 技术栈
- Vue 3 + Composition API + `<script setup>` 语法
- TypeScript strict 模式
- shadcn-vue 组件库
- Tailwind CSS

### 编辑规则
- **禁止整文件重写**：所有改动必须用最小 diff（Edit 工具的 old_string/new_string 精准替换）。除非用户明确要求重写或文件确实需要完全重构，否则不要用 Write 覆盖已有文件。
- 修改前先 Read 当前内容，保留用户已有代码、注释、格式。

### Vue 响应式陷阱
- `watchEffect` / `watch` 不得写成自我触发循环：禁止在副作用里直接修改它依赖的响应式状态；如必须，使用 `flush: 'post'`、guard 标志位或改用计算属性 (`computed`)。
- 改动 `watch`/`watchEffect` 后必须心算一遍依赖图：「这次写入会不会重新触发本回调？」
- `v-for` 的 `:key` 必须是稳定唯一 ID，**禁止用 index**（除非列表只读且永不重排）。
- 避免在模板中调用会产生新引用的函数（每次渲染都返回新对象/数组），用 `computed` 缓存。

### 验证（改完必须跑）
1. `pnpm vue-tsc --noEmit`（或项目里等价的类型检查脚本）必须 0 错误。
2. lint 必须通过（`pnpm lint` 或项目配置的命令）。
3. 改 UI 后在浏览器里实际验证一次，确认无闪烁、无控制台 warning。

未通过这三项前，不要声明任务完成。
