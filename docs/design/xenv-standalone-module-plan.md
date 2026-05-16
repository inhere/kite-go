# xenv 独立 Go Module 迁移方案

## 目标

将当前仓库中的 `xenv` 功能拆分为独立 Go module：

```text
module github.com/inhere/xenv
```

临时落地目录：

```text
./tmp/xenv
```

拆分后的项目应满足：

- 保留 `xenv` 相关的大部分提交记录。
- 独立提供 `xenv` CLI。
- 继续基于 `gookit/gcli`、`gookit/config`、`gookit/goutil` 等包。
- 不再依赖 `github.com/inhere/kite-go` 的 `internal` 或 `pkg/util`。
- 后续可以直接从 `./tmp/xenv` 拷贝到其他目录或推送到 `github.com/inhere/xenv`。

## 推荐策略

不要直接复制文件。直接复制会丢失提交历史。

推荐使用 `git filter-repo`：

1. 从当前仓库 clone 一个本地副本到 `tmp/xenv`。
2. 在 `tmp/xenv` 中只保留 `xenv` 相关路径。
3. 同时将部分路径重写成独立项目更合理的结构。
4. 在 `tmp/xenv` 中创建独立 `go.mod`。
5. 替换 import path。
6. 迁移少量 `pkg/util` 依赖到新项目内部。
7. 验证 `go test`、`go build`、`go run ./cmd/xenv --help`。

这种方式会重写 `tmp/xenv` 仓库历史，提交 hash 会变化，但提交作者、时间、message 和相关内容历史会保留。

## 建议目录结构

第一阶段建议低风险拆分，不急着把 `pkg/xenv` 提升到模块根目录，先保留原结构，减少 import 改动。

```text
tmp/xenv/
├── go.mod
├── LICENSE
├── README.md
├── cmd/
│   └── xenv/
│       └── main.go
├── pkg/
│   └── xenv/
│       ├── config/
│       ├── manager/
│       ├── models/
│       ├── service/
│       ├── shell/
│       ├── tools/
│       ├── xenvcom/
│       ├── xenvutil/
│       └── xenv.go
├── internal/
│   ├── xenvcmd/
│   │   ├── xenv_cmd.go
│   │   └── subcmd/
│   └── util/
├── data/
│   └── config.yaml
└── docs/
    ├── xenv-feature-report.md
    └── kite-xenv-spec-craft.md
```

关键设计：

- `internal/cli/xenvcmd` 迁移为新项目内的 `internal/xenvcmd`。
- 原因：拆分后 `xenvcmd` 只服务于本项目的 `cmd/xenv` 入口，不需要作为公共包暴露给其他 Go module。
- `pkg/xenv` 暂时保留，降低一次迁移的改动面。
- `internal/util` 承接当前 `github.com/inhere/kite-go/pkg/util` 中被 `xenv` 使用的少量工具函数。

## 需要保留的路径

建议从当前仓库筛选这些路径：

```text
cmd/xenv/
internal/cli/xenvcmd/
pkg/xenv/
data/xenv/
docs/xenv-feature-report.md
docs/feat-craft/kite-xenv-spec-craft.md
LICENSE
README.md
```

说明：

- `cmd/xenv` 是独立 CLI 入口。
- `internal/cli/xenvcmd` 是 CLI 命令定义，迁移后应改为 `internal/xenvcmd`。
- `pkg/xenv` 是核心功能模块。
- `data/xenv/config.yaml` 是配置样例。
- `docs/xenv-feature-report.md` 是当前功能报告。
- `docs/feat-craft/kite-xenv-spec-craft.md` 是原始需求草稿。
- `LICENSE` 和 `README.md` 可作为独立项目基础文件。

## 历史保留拆分命令

先确认当前仓库状态：

```powershell
git status --short
```

建议把临时目录加入父仓库本地 exclude，避免 `tmp/xenv` 出现在父仓库状态里：

```powershell
Add-Content .git/info/exclude "`n/tmp/"
```

如果未安装 `git-filter-repo`：

```powershell
python -m pip install git-filter-repo
```

创建独立临时仓库：

```powershell
New-Item -ItemType Directory -Force tmp | Out-Null
git clone --no-hardlinks . tmp/xenv
Set-Location tmp/xenv
```

筛选历史并重命名路径：

```powershell
git filter-repo --force `
  --path cmd/xenv/ `
  --path internal/cli/xenvcmd/ `
  --path pkg/xenv/ `
  --path data/xenv/ `
  --path docs/xenv-feature-report.md `
  --path docs/feat-craft/kite-xenv-spec-craft.md `
  --path LICENSE `
  --path README.md `
  --path-rename internal/cli/xenvcmd/:internal/xenvcmd/ `
  --path-rename data/xenv/:data/ `
  --path-rename docs/feat-craft/kite-xenv-spec-craft.md:docs/kite-xenv-spec-craft.md
```

检查历史：

```powershell
git log --oneline --all --decorate -30
git log --oneline -- pkg/xenv cmd/xenv internal/xenvcmd | Select-Object -First 30
```

预期：

- 只剩 `xenv` 相关提交。
- 之前大量 `xenv` 提交仍在。
- 不相关的 Kite 功能提交会被剪掉。
- 空提交会被移除。

## go.mod 设计

在 `tmp/xenv/go.mod` 中创建：

```go
module github.com/inhere/xenv

go 1.24

toolchain go1.24.6

require (
	github.com/BurntSushi/toml v1.6.0
	github.com/goccy/go-json v0.10.6
	github.com/gookit/cliui v0.2.3
	github.com/gookit/config/v2 v2.2.8
	github.com/gookit/gcli/v3 v3.3.1
	github.com/gookit/goutil v0.7.5
)
```

不建议将原仓库中的本地 replace 带入独立模块：

```go
replace github.com/gookit/gcli/v3 => ../gcli
```

除非本地明确需要联调 `gcli`，否则独立项目不应默认依赖本地相对路径。

整理依赖：

```powershell
go mod tidy
```

## Import 迁移规则

需要替换：

```text
github.com/inhere/kite-go/pkg/xenv
=> github.com/inhere/xenv/pkg/xenv
```

```text
github.com/inhere/kite-go/pkg/xenv/...
=> github.com/inhere/xenv/pkg/xenv/...
```

```text
github.com/inhere/kite-go/internal/cli/xenvcmd
=> github.com/inhere/xenv/internal/xenvcmd
```

当前还有一个重要耦合：

```text
github.com/inhere/kite-go/pkg/util
```

这些调用主要包括：

```text
NormalizePath
SplitPath
JoinPaths
EnsureDir
CopyFile
CreateSymlink
```

建议在新模块创建：

```text
internal/util/path.go
internal/util/file.go
```

然后替换 import：

```go
"github.com/inhere/kite-go/pkg/util"
```

为：

```go
"github.com/inhere/xenv/internal/util"
```

检查是否还有旧路径：

```powershell
rg "github.com/inhere/kite-go"
```

目标是没有任何输出。

## 迁移任务拆分

### 任务 1：生成带历史的临时仓库

执行：

```powershell
New-Item -ItemType Directory -Force tmp | Out-Null
git clone --no-hardlinks . tmp/xenv
Set-Location tmp/xenv
git filter-repo --force `
  --path cmd/xenv/ `
  --path internal/cli/xenvcmd/ `
  --path pkg/xenv/ `
  --path data/xenv/ `
  --path docs/xenv-feature-report.md `
  --path docs/feat-craft/kite-xenv-spec-craft.md `
  --path LICENSE `
  --path README.md `
  --path-rename internal/cli/xenvcmd/:internal/xenvcmd/ `
  --path-rename data/xenv/:data/ `
  --path-rename docs/feat-craft/kite-xenv-spec-craft.md:docs/kite-xenv-spec-craft.md
```

验证：

```powershell
git log --oneline --all --decorate -30
Get-ChildItem
```

提交：此步骤是过滤仓库历史，不需要在父仓库提交。

### 任务 2：初始化独立 go.mod

创建 `go.mod`：

```go
module github.com/inhere/xenv

go 1.24

toolchain go1.24.6

require (
	github.com/BurntSushi/toml v1.6.0
	github.com/goccy/go-json v0.10.6
	github.com/gookit/cliui v0.2.3
	github.com/gookit/config/v2 v2.2.8
	github.com/gookit/gcli/v3 v3.3.1
	github.com/gookit/goutil v0.7.5
)
```

执行：

```powershell
go mod tidy
```

提交：

```powershell
git add go.mod go.sum
git commit -m "chore: initialize xenv module"
```

### 任务 3：更新命令包路径

确认路径已在 filter 阶段重命名：

```text
internal/xenvcmd/
```

更新 `cmd/xenv/main.go`：

```go
import (
	"github.com/inhere/xenv/internal/xenvcmd"
	"github.com/inhere/xenv/pkg/xenv/xenvcom"
)
```

提交：

```powershell
gofmt -w cmd pkg
git add cmd pkg
git commit -m "refactor: move xenv command package"
```

### 任务 4：替换 xenv 内部 import path

批量替换：

```text
github.com/inhere/kite-go/pkg/xenv
=> github.com/inhere/xenv/pkg/xenv
```

检查：

```powershell
rg "github.com/inhere/kite-go/pkg/xenv"
```

格式化并整理依赖：

```powershell
gofmt -w cmd pkg
go mod tidy
```

提交：

```powershell
git add .
git commit -m "refactor: update imports for xenv module"
```

### 任务 5：迁移 util 依赖

创建：

```text
internal/util/path.go
internal/util/file.go
```

最小函数集：

```text
NormalizePath
SplitPath
JoinPaths
EnsureDir
CopyFile
CreateSymlink
```

替换：

```go
"github.com/inhere/kite-go/pkg/util"
```

为：

```go
"github.com/inhere/xenv/internal/util"
```

检查：

```powershell
rg "github.com/inhere/kite-go/pkg/util"
```

提交：

```powershell
gofmt -w internal pkg
go mod tidy
git add internal pkg go.mod go.sum
git commit -m "refactor: move local utilities into xenv module"
```

### 任务 6：整理文档和 README

建议保留：

```text
data/config.yaml
docs/xenv-feature-report.md
docs/kite-xenv-spec-craft.md
```

新增或改写 `README.md`，至少包含：

```markdown
# xenv

Local development environment and SDK manager.

## Install

```bash
go install github.com/inhere/xenv/cmd/xenv@latest
```

## Quick Start

```bash
xenv init
eval "$(xenv shell --type bash)"
xenv tools index
xenv tools list
xenv use go:latest
```
```

提交：

```powershell
git add README.md docs data
git commit -m "docs: add standalone xenv usage"
```

### 任务 7：验证独立模块

基础验证：

```powershell
go mod tidy
go test ./pkg/xenv/service ./pkg/xenv/shell ./pkg/xenv/manager ./internal/xenvcmd/...
go build ./cmd/xenv
go run ./cmd/xenv --help
go run ./cmd/xenv shell --type pwsh
```

完整验证：

```powershell
go test ./...
```

注意：当前原项目里 `pkg/xenv/tools` 的测试和实现已有不一致：

```text
ParseVersionSpec("go")
```

当前实现会解析为：

```text
go:latest
```

但测试预期报错。拆分时需要明确规则：

- 如果希望 `xenv use go` 表示 `go:latest`，就改测试。
- 如果希望必须写 `go:version`，就改实现。

建议保留当前 CLI 体验，即：

```text
xenv use go == xenv use go:latest
```

然后单独提交：

```powershell
git commit -m "test: align version spec behavior with cli defaults"
```

## 与 kite-go 的后续关系

拆分完成后，`kite-go` 不再内嵌 `xenv` 源码，也不再通过 Go module 依赖 `github.com/inhere/xenv` 注册 `xenvcmd`。

目标关系：

- `kite-go` 删除 `cmd/xenv`。
- `kite-go` 删除 `pkg/xenv`。
- `kite-go` 删除 `internal/cli/xenvcmd`。
- `kite` 不再内置提供 `kite xenv`。
- 用户直接安装独立 CLI：`github.com/inhere/xenv/cmd/xenv`。
- 如需在 `kite` 中调用 `xenv`，通过 Kite 的 ext/外部命令注册机制调用系统里的 `xenv` 可执行文件。

推荐安装方式：

```bash
go install github.com/inhere/xenv/cmd/xenv@latest
```

安装后直接使用：

```bash
xenv init
xenv shell --type bash
xenv tools list
xenv use go:latest
```

如果希望仍然通过 `kite` 入口调用，使用 ext 方式注册外部命令，例如将 `xenv` 注册为一个转发到系统命令的扩展：

```text
kite ext -> xenv executable
```

目标是让 `kite-go` 只负责发现和调用外部 `xenv`，不再编译、链接或维护 `xenv` 代码。

迁移完成后，`kite-go` 侧需要做一次清理提交：

```text
remove cmd/xenv
remove pkg/xenv
remove internal/cli/xenvcmd
remove xenv command registration from internal/cli/boot.go
update docs to point users to github.com/inhere/xenv
optionally add ext registration example for xenv
```

## 推荐执行路线

建议分两阶段。

### 阶段 1：`tmp/xenv` 独立可构建

- 使用 `git filter-repo` 生成带历史的新仓库。
- 保留 `pkg/xenv` 结构。
- 将 `internal/cli/xenvcmd` 改为新项目内的 `internal/xenvcmd`。
- 复制最小 `internal/util`，断开对 `kite-go/pkg/util` 的依赖。
- 确保 `go build ./cmd/xenv` 通过。

### 阶段 2：从 `kite-go` 移除内嵌 xenv

- 从 `kite-go` 删除 `cmd/xenv`、`pkg/xenv`、`internal/cli/xenvcmd`。
- 从 `internal/cli/boot.go` 移除 `xenvcmd.XEnvCmd` 注册。
- 删除 `kite-go` 中仅为 `xenv` 保留的构建目标，例如 `build-xenv`、`install-xenv`。
- 更新 `kite-go` 文档，说明 `xenv` 已迁移到 `github.com/inhere/xenv`。
- 如需保留 Kite 内调用体验，补充通过 ext 注册外部 `xenv` CLI 的示例。

## 风险点

主要风险：

- `git filter-repo` 会重写 `tmp/xenv` 仓库历史，但不应在原仓库直接运行。
- `internal/cli/xenvcmd` 需要从原路径迁移到新项目的 `internal/xenvcmd`，避免保留 `internal/cli` 这种来自 `kite-go` 的目录语义。
- `pkg/util` 依赖必须处理，否则新 module 仍会绑定 `kite-go`。
- 版本解析测试当前已有不一致，拆分后需要作为单独决策修正。
- 父仓库最好 exclude `tmp/`，避免把临时独立仓库误提交进 `kite-go`。

## 完成标准

`./tmp/xenv` 满足以下条件即认为第一阶段完成：

```powershell
rg "github.com/inhere/kite-go"
```

无输出。

```powershell
go mod tidy
go build ./cmd/xenv
go run ./cmd/xenv --help
go run ./cmd/xenv shell --type pwsh
```

全部通过。

基础测试通过：

```powershell
go test ./pkg/xenv/service ./pkg/xenv/shell ./pkg/xenv/manager ./internal/xenvcmd/...
```

历史检查通过：

```powershell
git log --oneline -- pkg/xenv cmd/xenv internal/xenvcmd | Select-Object -First 30
```

能看到原来 `xenv` 相关提交记录。
