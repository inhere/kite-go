# `kite run` 命令运行流程分析

> 分析日期：2026-02-28

## 概述

`kite run`（别名 `kite exec`）是一个通用命令执行器，支持运行 kite 命令别名、kite 扩展、脚本任务、脚本文件和系统命令。

## 相关源码文件

| 文件 | 职责 |
|------|------|
| `internal/cli/toolcmd/runany_cmd.go` | CLI 命令定义与入口函数 |
| `internal/biz/cmdbiz/runany.go` | 核心业务逻辑：智能分发 RunAny() |
| `pkg/kscript/runner.go` | Runner 结构体定义与脚本加载逻辑 |
| `pkg/kscript/runner_run.go` | 脚本任务/脚本文件的实际执行逻辑 |

---

## 1. 命令解析入口

**入口函数：** `runAnything()` — `internal/cli/toolcmd/runany_cmd.go:85`

```
kite run [OPTIONS] <name> [args...]
```

### 支持选项

| 选项 | 说明 |
|------|------|
| `--type, -t` | 强制指定运行类型：alias / script / ext / plugin / system |
| `--shell` | 指定 shell 包装类型：bash / sh / zsh / cmd / pwsh 等 |
| `--vars, --var` | 自定义变量，格式 `name=value` |
| `--env, -e` | 自定义环境变量，格式 `KEY=VALUE` |
| `--chdir, --cd` | 自动向上查找目录并切换到该目录作为工作目录 |
| `--list, -l` | 列出所有脚本或别名信息 |
| `--show, -i` | 显示指定名称的详细信息 |
| `--search, -s` | 按名称搜索匹配的脚本 |
| `--verbose, --verb` | 输出执行上下文详细信息 |
| `--dry-run` | 仅打印命令，不实际执行 |

### 入口分支逻辑

```
kite run <name> [args]
 │
 ├─ --list / -l          → listInfos()              列出所有别名/脚本
 ├─ --show / -i          → showInfo(name)            显示详细信息
 ├─ --type=system        → cmdr.NewCmd().FlushRun()  直接执行系统命令
 ├─ --type=alias         → RunKiteCmdByAlias()       执行 kite 命令别名
 ├─ --type=script        → app.Scripts.Run()         直接执行脚本（跳过别名/ext检查）
 └─ (默认)               → cmdbiz.RunAny()           智能分发（见第2节）
```

---

## 2. 核心分发：`RunAny()`

**位置：** `internal/biz/cmdbiz/runany.go:39`

按以下优先级依次尝试，首次匹配即执行：

```
name
 │
 ├─①  app.Kas.HasAlias(name)
 │      → RunKiteCmdByAlias()
 │        解析别名字符串 → 拆分 cmd + args → app.Cli.RunCmd()
 │
 ├─②  app.Exts.Exists(name)
 │      → app.Exts.Run(name, &kiteext.RunCtx{Args: args})
 │
 ├─③  app.Scripts.TryRun(name, args, ctx)
 │      → kscript.Runner（见第3节）
 │
 ├─④  (TODO) plugin
 │
 └─⑤  sysutil.HasExecutable(name)
        → cmdr.NewCmd(name, args...).FlushRun()
        未找到 → 返回错误
```

---

## 3. 脚本加载：`Runner.InitLoad()`

**位置：** `pkg/kscript/runner.go:127`

采用懒加载策略，只在首次调用时执行，后续调用直接跳过。

```
InitLoad()
 ├─ LoadScriptTasks()     加载脚本任务定义
 │    ├─ 遍历 DefineFiles（支持 ? 前缀表示可选文件，支持 $os/$user 变量）
 │    ├─ 支持格式：yaml / toml / ini
 │    ├─ findAutoTaskFiles()
 │    │    从当前目录向上遍历（最多 AutoMaxDepth 层）
 │    │    查找 AutoTaskFiles 中配置的文件名 + AutoTaskExts 中配置的扩展名
 │    │    倒序（顶层优先）加载
 │    └─ 解析 __settings 键 → 加载到 taskSettings
 │
 ├─ LoadScriptApps()      加载脚本 app 定义
 │    从 ScriptAppDirs 读取定义文件，存入 appFiles 映射
 │
 └─ LoadScriptFiles()     加载脚本文件
      从 ScriptDirs 扫描所有文件，存入 scriptFiles 映射 {filename: fullpath}
```

---

## 4. 脚本执行：`Runner.TryRun()`

**位置：** `pkg/kscript/runner_run.go:35`

```
TryRun(name, args, ctx)
 │
 ├─ InitLoad()                       懒加载（见第3节）
 │
 ├─ name 含空格？
 │    → SearchByKeywords()           模糊关键词匹配
 │       匹配唯一 → 替换 name
 │       匹配多个 → 返回错误提示
 │       无匹配   → 返回错误
 │
 ├─ LoadScriptTaskInfo(name)
 │    在 Scripts map 中查找 → parseScriptTask() 解析为 ScriptTask
 │    找到 → runScriptTask()           ─── 见第4.1节
 │
 └─ LoadScriptFileInfo(name)
      在 scriptFiles 中查找（支持带扩展名或自动匹配扩展名）
      找到 → runScriptFile()           ─── 见第4.2节
      均未找到 → 返回 (found=false, nil)
```

---

### 4.1 执行脚本任务：`runScriptTask()`

**位置：** `pkg/kscript/runner_run.go:110`

```
runScriptTask(st *ScriptTask, inArgs, ctx)
 │
 ├─ BeforeFn()                        可选前置钩子（verbose 模式下打印脚本信息）
 ├─ 校验必要参数数量 st.ParseArgs()
 │
 ├─ buildTaskRenderVars()             构建模板变量 map：
 │    - st.Vars（任务级变量，支持动态变量 resolveDynVars）
 │    - ctx.Vars（--vars 传入的变量，优先级更高）
 │    - time.*（内置时间变量，含 unix_sec / datetime / date_ymd 等）
 │    - vars / groups（来自 __settings 全局变量组）
 │    - gvs / paths / kite（由 AppendVarsFn 注入的应用全局变量）
 │    - cur_dir（当前工作目录）
 │
 ├─ MergeEnv()                        合并 ENV：settings.Env → st.Env → ctx.Env
 ├─ ParseVarInEnv()                   在 ENV 值中渲染变量引用
 ├─ renderTaskVars(workdir, ...)      渲染 workdir 中的变量
 ├─ AppendArgsToVars()                注入参数变量：$0/$1/.../$@/$*
 │
 ├─ 执行 Deps 依赖任务（递归调用 runScriptTask）
 │
 └─ 遍历 st.Cmds 依次执行每条命令：
      ├─ tc.isRef → 引用其他任务（递归执行）
      ├─ tc.appendVars(vars)           加载命令独有变量
      ├─ renderTaskVars(tc.Run, vars)  渲染命令字符串中的变量
      │    支持：$var, ${var}, ${var|default}, 链式访问, ENV 变量
      └─ 构建并执行命令：
           shell != "" → cmdr.NewCmd(shell, "-c", renderedLine)
           shell == "" → cmdr.NewCmdline(renderedLine)
           .WorkDirOnNE(cmdDir)
           .WithDryRun(ctx.DryRun)
           .AppendEnv(envMap)
           .PrintCmdline2()            打印实际执行命令（含参数）
           .FlushRun()                 执行并实时输出 stdout/stderr
```

---

### 4.2 执行脚本文件：`runScriptFile()`

**位置：** `pkg/kscript/runner_run.go:327`

```
runScriptFile(sf *ScriptFile, inArgs, ctx)
 └─ cmdr.NewCmd(sf.BinName, sf.File)
      .WorkDirOnNE(sf.Workdir)
      .WithDryRun(ctx.DryRun)
      .AppendEnv(sf.Env)
      .AddArgs(inArgs)
      .PrintCmdline2()
      .FlushRun()
```

**解释器解析规则（`LoadScriptFileInfo`）：**

1. 文件名包含扩展名 → 直接在 `scriptFiles` 中查找
2. 文件名不含扩展名 → 遍历 `AllowedExt` 自动追加扩展名查找
3. 解释器（BinName）：默认取扩展名去掉点（`.sh` → `sh`），可由 `ExtToBinMap` 配置覆盖（如 `.sh` → `bash`）

---

## 5. 完整流程总图

```
kite run <name> [args]
 └── runAnything()                          # toolcmd/runany_cmd.go
      │
      ├── [--type=system]  cmdr.FlushRun()  # 直接运行系统命令
      ├── [--type=alias]   RunKiteCmdByAlias
      ├── [--type=script]  Scripts.Run()
      │
      └── cmdbiz.RunAny()                   # cmdbiz/runany.go
           │
           ├── ① kite 命令别名 → RunKiteCmdByAlias → app.Cli.RunCmd()
           ├── ② kite ext     → app.Exts.Run()
           ├── ③ kscript      → Runner.TryRun()     # kscript/runner_run.go
           │    ├── InitLoad()（懒加载脚本定义）
           │    ├── ScriptTask → runScriptTask()
           │    │    ├── 构建变量 → 渲染命令 → 执行 deps → 遍历执行 cmds
           │    │    └── 支持：shell 包装 / 变量模板 / dry-run / verbose
           │    └── ScriptFile → runScriptFile()
           │         └── 按扩展名选择解释器执行文件
           └── ⑤ 系统命令   → cmdr.NewCmd().FlushRun()
```

---

## 6. 关键设计点

- **懒加载**：脚本定义文件在首次需要时才加载，并缓存结果（`taskLoaded` / `fileLoaded` 标志）
- **自动目录搜索**：`AutoTaskFiles` 支持从当前目录向上逐层查找任务配置文件（最多 `AutoMaxDepth` 层）
- **变量模板渲染**：使用 `textutil.NewStrVarRenderer()` 支持 PHP/Shell 风格变量语法，含链式访问和默认值
- **依赖任务**：`st.Deps` 支持任务前置依赖，递归执行
- **任务引用**：`tc.isRef` 支持在命令列表中引用其他任务（类似函数调用）
- **Dry Run**：所有执行路径均支持 `--dry-run`，只打印不执行
