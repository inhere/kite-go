# kscript 可嵌入 Go 库调研（2026-09-14）

## 背景与结论

`kite-go/pkg/kscript` 当前面向“任务定义 + 脚本文件 + 变量/env + shell 执行 + 自动发现”，并支持 YAML/JSON/TOML。它更像一个轻量任务编排层，而不是单一脚本语言运行时。GitHub 上没有一个库可以无缝替代全部能力；建议按需求组合：

* 若首要目标是兼容 Taskfile 风格的跨平台任务执行，优先评估 **go-task/task**，但通常通过其 CLI/内部包接入，API 稳定性和嵌入边界需单独确认。
* 若必须在进程内执行 shell、包括 Windows 无系统 shell 场景，使用 **mvdan/sh**；它提供 parser/formatter/interpreter，可自行注册命令与限制访问。
* 若脚本需要通用语言和丰富生态，选 **goja**（JavaScript ES5.1）；若要更小、更可控的 DSL，选 **Tengo** 或 **Starlark**。
* 仅需条件、映射和校验表达式时，选 **expr**，不要把它当作任务执行器。

## 项目对比

| 项目 | 官方定位/API 证据 | 许可证 / Stars* | 对 kscript 的适配性 | 主要限制与集成工作 |
|---|---|---:|---|---|
| [go-task/task](https://github.com/go-task/task) · [Taskfile 文档](https://taskfile.dev/) | README 定位为“fast, cross-platform build tool inspired by Make”；Taskfile YAML、变量、依赖、并发、平台条件等由 CLI 实现。 | MIT · 16,141 | 与现有 YAML 任务模型最接近，可复用 Taskfile 生态和 shell 命令约定。适合把 kscript 的发现/变量层映射到 Taskfile。 | 项目核心是 CLI，公开可嵌入 API 不是主要承诺；需评估导入内部包或作为子进程调用。Task 仍依赖外部 shell/命令，沙箱能力有限。 |
| [mvdan/sh](https://github.com/mvdan/sh) · [shell pkg](https://pkg.go.dev/mvdan.cc/sh/v3/shell) | README 明确是 shell parser、formatter、interpreter，支持 POSIX Shell/Bash/Zsh/mksh；`mvdan.cc/sh/v3/shell` 提供带 shell 语义的调用辅助；解释器可在无系统 shell 下运行并通过 handlers 沙箱命令和访问。 | BSD-3-Clause · 9,056 | 最适合作为 kscript 的进程内 shell 后端：保留脚本文件、变量和跨平台执行，同时可控制命令集合、文件/网络访问。 | 只解决 shell 语言执行，不提供 YAML 任务 DAG、自动发现、重试/缓存等编排；需自行实现上下文、超时、日志和错误模型。README 要求 Go 1.26+（以仓库当前版本为准）。 |
| [goja](https://github.com/dop251/goja) | README：纯 Go 的 ECMAScript 5.1(+ ) 实现，强调标准兼容；可创建 runtime、导出 Go 值/函数、执行 JS。 | MIT · 7,089 | 可把每个任务脚本作为 JS，向 runtime 注入 `run`、变量和 Go 服务；语言表达力和库生态优于自定义 DSL。 | ES5.1，不是现代 Node.js；默认没有 Node 内置模块。执行隔离、CPU/内存配额和阻止恶意代码需由宿主设计；JS 与 Go 频繁复杂数据交互有额外成本。 |
| [Tengo](https://github.com/d5/tengo) | README 介绍为“fast, secure, embeddable scripting language for Go”，编译为栈式 VM bytecode；提供 `Eval` 等 API，并可注册 Go 函数/变量。 | MIT · 3,836 | 适合内嵌、可控的任务 DSL：脚本小、启动快，可将 kscript 的变量和操作封装成内置函数。 | 生态和现成库少于 JS；不负责 shell 进程、YAML 任务图或文件发现。需要定义标准库、错误/超时和权限边界；检查项目版本与 Go 兼容性。 |
| [google/starlark-go](https://github.com/google/starlark-go) · [Go API](https://pkg.go.dev/go.starlark.net/starlark) | 官方 README：Starlark interpreter in Go，导入路径 `go.starlark.net/starlark`；Starlark 是面向配置的小型 Python 方言，线程可并行执行。 | BSD-3-Clause · 2,760 | 适合声明式任务参数、条件和宏；语言刻意限制副作用，宿主可通过自定义 builtins 暴露安全操作。 | 不是 shell，也没有内置 I/O/进程；必须自己提供 `run` 等 builtins 和任务调度。Python 语法仅为方言，不能直接运行普通 Python 包。 |
| [expr](https://github.com/expr-lang/expr) · [Language/Guide](https://expr-lang.org/) | README 定位为 Go-centric、安全、无副作用表达式求值器；优化编译器 + bytecode VM，输入驱动。 | MIT · 8,011 | 适合替换 kscript 中的 `if`、参数模板、过滤和校验表达式；可编译后缓存，类型检查较好。 | 明确是 side-effect-free expression，不执行脚本、命令或任务 DAG；需与 shell/任务运行器组合。表达式语法和可用函数集需在 API 层固定，避免把业务逻辑塞进字符串。 |

\* Stars 和许可证为 GitHub REST API 在 2026-09-14 UTC 查询值，会随时间变化；许可证字段以仓库元数据为准。

## 对 kscript 的落地建议

1. **先抽象执行接口**：将当前 runner 中“解析任务、选择脚本、执行命令、收集结果”拆为 `TaskLoader`、`ScriptEngine`、`CommandRunner`。这样可以先接入 mvdan/sh，而不改变 YAML/JSON/TOML 和自动发现。
2. **表达式单独接入**：把条件和模板求值限定为 expr；禁止表达式直接产生任意进程调用。若需要脚本逻辑，再增加 Tengo/Starlark engine，并通过显式 builtin 暴露能力。
3. **JS 作为可选扩展**：goja 适合用户确实需要 JavaScript 的场景，建议独立 build tag/模块，设置超时、上下文取消和可注入 API 白名单。
4. **Taskfile 兼容路线**：若希望复用 Taskfile 文件和生态，先做格式/语义映射验证，再决定调用 Task CLI（进程隔离、升级简单）还是导入内部包（避免 CLI 开销但承担 API 变动）。不要在未验证前把 go-task 当作稳定嵌入 SDK。
5. **安全与可观测性**：无论选哪种 engine，都要统一 context 超时/取消、stdout/stderr、退出码、结构化错误和审计日志；脚本来源不可信时，进程、文件、网络和资源配额必须由宿主层控制。

## 参考的一手来源

* GitHub 仓库 README、源码和 LICENSE：上述各仓库链接。
* Go 包文档：[mvdan.cc/sh/v3/shell](https://pkg.go.dev/mvdan.cc/sh/v3/shell)、[go.starlark.net/starlark](https://pkg.go.dev/go.starlark.net/starlark)、[github.com/d5/tengo/v2](https://pkg.go.dev/github.com/d5/tengo/v2)。
* GitHub 元数据 API（查询时间 2026-09-14）：`https://api.github.com/repos/{owner}/{repo}`。
