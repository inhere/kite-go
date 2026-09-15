<!-- template_id: plan; template_version: 1.2.0 -->
# github.com/gookit/kscript 独立库实施计划

> 状态：Draft 0.2 / 待人工计划批准
>
> 目标 module：`github.com/gookit/kscript`；Git root：`D:/work/inhere/my-tools-dev/gookit2/kscript`。
> 该目录当前尚未创建；本计划不把路径确认视为已授权初始化 Git 或创建远端仓库。
>
> 下文 `<target>` 均指上述 `D:/work/inhere/my-tools-dev/gookit2/kscript`。

## 修订记录

| 版本 | 日期 | 作者 | 摘要 |
|---|---|---|---|
| 0.1 | 2026-09-15 | Codex | 将 Draft 0.2 设计拆为基线、独立核心、格式/执行器、Kite 迁移、第二应用验证和发布准备波次 |
| 0.2 | 2026-09-15 | Codex | 将新库 Git root 明确为 `gookit2/kscript`，同步目标路径、前置 Gate、任务和回滚边界 |

仅任务合同、范围、风险、验证、生命周期或执行语义变化时递增版本；路径、作者、状态等元数据纠正不增加版本。

## 目标与完成定义

将现有 `kite-go/pkg/kscript` 抽离为 `github.com/gookit/kscript`，让外部 Go 1.23+ 应用能够在不初始化 Kite 的情况下加载定义、检查计划、运行任务和脚本文件，并取得隔离、可取消、可分类的结构化结果。Kite 继续通过适配层保留 alias、extension、plugin、系统命令兜底和旧配置语义。

完成定义：

1. 独立 module 的主包和 `formats` 包可在外部临时 module 中以 Go 1.23 编译运行；不导入 `github.com/inhere/kite-go`、CLI 框架或全局 Kite 状态。
2. Definition、Task、Step、Request、Runner、Result、Engine、Handler 的已批准 API 有 GoDoc、错误契约和兼容测试。
3. YAML/JSON/TOML 经同一 schema 版本 1 模型解码；未定义字段、重复名称、缺失引用和循环在副作用前报错。
4. Task/Step 条件判断使用 expr，真实 Run 与 Inspect 的求值边界符合设计；错误、Skipped、Deferred 可观察。
5. 顺序 deps/task call、参数/变量/env/目录优先级、Shell/file/exec/host、dry-run、取消/超时、输出上限和错误分类有测试证据。
6. Kite legacy converter 通过现有任务 fixture；第二个真实 Go 应用的项目路径、Go 版本和实际运行证据已记录。
7. 所有阶段均有 focused validation、无越界 dirty work、原子本地 commit；发布、推送、tag 和部署另行授权。

## 范围、排除项与授权

范围冻结为设计 Draft 0.2 与评审 PASS 所覆盖的 v0.1：独立 module、类型模型、顺序调度、条件判断、外部进程 Engine、脚本文件、Handler、格式 loader、Kite 适配和两个 consumer 验证。

排除项：有界并行、重试、缓存、完整 Taskfile/justfile 兼容、Just 运行时依赖、Goja/Tengo/Starlark VM、远程执行、分布式调度、cron、容器编排、发布部署和业务补偿回滚。新增外部服务、数据库、硬件、生产流量或外部消息不在本计划内。

- `host_or_non_offline_action=NOT_APPLICABLE`

本计划只描述本地源码、文档和离线/本机测试。创建新 Git 仓库、推送、发布、外部应用部署或真实生产动作均须在相关任务前重新声明并取得外部动作 Gate；设计批准和计划批准均不自动授权这些动作。

## 输入与批准证据

- 设计：[2026-09-15-kscript-library-design.md](../design/2026-09-15-kscript-library-design.md)，revision 0.2，candidate commit `bf307538cdbed3f9a0a1eafa40784a45256db7b3`。
- 评审：[2026-09-15-kscript-library-design-review.md](../review/2026-09-15-kscript-library-design-review.md)，针对同一 candidate 的 Standards/Governance 与 Spec/Executability 双轴结论 PASS，未发现核心阻断。
- 当前用户请求：“OK 审核一遍，没问题出实施计划”；该请求授权生成本计划，不等同于批准实施。
- 工作区绑定：`D:/work/inhere/my-tools-dev/standards.json`，IDEV-STD 0.19.0、profile `go-tools`、revision `9dc4a010355e1eb7a5c53f99558495df09eb9553`；`probe=BOUND`、`validate=PASS`。

## Capability Discovery

### Capability decisions

| capability_id | required_capability | searched_candidates | direct_reuse | thin_adapter_or_owner_extension | decision | proven_gap | duplication_and_lifecycle_risk |
|---|---|---|---|---|---|---|---|
| CAP-01 | 独立 Go module 与公共任务 API | 现有 `pkg/kscript`；go-task/task；just | 现有 kscript 模型可迁移，但 Kite module 不能作为公共依赖 | 从现有模型迁移到新主包，Kite 通过 converter | MINIMAL_NEW_MODULE | 需要独立 import path、无 Kite internal 依赖和稳定错误/API | 新 module owner 为 gookit/kscript；Kite converter 只由 kite-go 维护 |
| CAP-02 | 顺序依赖、循环检查和运行隔离 | 现有 `deps`/`@task:` 递归；go-task/task DAG | 复用现有顺序行为和测试 fixture | `internal/graph` 做引用校验与 CallState | OWNER_EXTENSION | 现有递归没有完整循环/膨胀保护，RunCtx 可变共享 | 调度器只存在新库；Kite 不复制 DAG |
| CAP-03 | exec/Shell/file/host 执行 | 现有 `cmdr`；stdlib os/exec；mvdan/sh | 首期复用标准进程能力；现有 cmdr 仅作为实现参考 | `Engine` 接口；Kite Handler adapter | THIN_ADAPTER | 需 argv/Shell 结构化参数、取消、进程清理、输出模型 | Engine owner 为新库；mvdan/sh 延期为可选实现 |
| CAP-04 | 任务/Step 条件判断和表达式 | 现有 expr 依赖；expr 官方库 | 直接复用 expr 编译器依赖 | 封装只读 Env、bool 校验和 Deferred Inspect | DIRECT_REUSE | 现有 resolveIfExpr panic、打印且固定 true，未接入 Step | expr 由新库统一维护；不得再造表达式语言 |
| CAP-05 | YAML/JSON/TOML loader 与发现 | `gookit/config/v2`；stdlib encoding/json；goccy YAML；现有 loader | 复用现有格式行为作为 fixture | `formats` 解码到统一 Definition；发现策略显式化 | THIN_ADAPTER | 现有 loader 状态和自动发现隐式，不能成为公共默认 | loader 生命周期归新库；Kite legacy discovery 留在 converter |
| CAP-06 | Go consumer 示例与验证 | 现有 Kite；独立临时 module；第二应用 | 无现成第二应用证据 | examples/basic 与 consumer fixture | OWNER_EXTENSION | 必须证明不导入 Kite 且 Go 1.23 可用 | 示例随新库维护；真实第二应用由其 owner 维护 |

### New module candidates

| candidate_id | capability_id | proposed_module | searched_candidates | direct_reuse_gap | thin_adapter_or_owner_extension_gap | proven_gap | unique_owner_and_lifecycle | deletion_or_merge_handling |
|---|---|---|---|---|---|---|---|---|
| MOD-01 | CAP-01 | `github.com/gookit/kscript` | 现有 kscript、go-task/task、just | Kite module 无法提供独立公共边界；task/just 定位或语言不符 | 新库需拥有 API、版本、测试和发布生命周期 | 设计与评审已证明独立 import/多应用目标 | gookit/kscript owner；v0.x → v1 SemVer | 若取消抽离，删除新 module，Kite 保留原包；不同时维护两个核心实现 |

### Rejected new tools

| rejected_candidate | capability_id | deletion_test_and_reason |
|---|---|---|
| go-task/task 直接作为核心 | CAP-01 | 删除它后仍可由新库模型和标准进程 Engine 满足独立 API；其嵌入 API 不是主要承诺，故只借鉴语义 |
| go-task/task 直接作为调度器 | CAP-02 | 删除它后仍可由新库自己的顺序图校验满足 v0.1；需要的调度边界不依赖其内部 API |
| just 作为 Go 依赖 | CAP-01 | 删除后仍可满足所有 v0.1 验收；just 是 Rust CLI，改为外部依赖会破坏 Go 原生嵌入目标 |
| just 作为执行后端 | CAP-03 | 删除后仍可由标准进程 Engine 满足 v0.1；just 的运行时不提供 Go 原生嵌入接口 |
| mvdan/sh/Tengo/Starlark/goja 首期内置 | CAP-03 | 删除后外部进程 + expr 仍满足 v0.1；语言 VM/进程内 Shell 属于后续能力，按最小新模块原则延期 |
| 语言 VM 首期内置 | CAP-04 | 删除后 expr 仍满足条件判断验收；Goja/Tengo/Starlark 的完整脚本语义属于后续设计 |

## 前置检查与 fail-closed 条件

1. 实施前确认目标 module Git root、目录 owner、第二应用路径和 Go 1.23 CI runner；未确认则停在人工 Gate，不在计划中猜测。
2. 记录目标 Git root、branch、HEAD、`git status --short`；现有 Kite baseline 为 root `D:/work/inhere/my-tools-dev/inhere-tools/kite-go`、branch `main`、HEAD `bf307538cdbed3f9a0a1eafa40784a45256db7b3`。
3. Kite 当前 dirty/untracked 为 `.codebase-memory/`、`docs/research/`；它们不是本计划所有权，必须保留，不得 add、删除或 reset。设计和评审文件为本计划的受控路径，已在 candidate commit 中固定。
4. 重新运行工作区 `probe`/`validate` 和目标 release 的 loader；revision、profile、root role 或 dirty 状态变化时停下并重新确认。
5. 目标路径与 package import path 不一致、发现已有 `github.com/gookit/kscript` 远端内容、或已有实现与设计不同，归 `Semantic Amendment`，先回设计评审。
6. 任何新 CLI、schema、runner、registry、adapter 或 module 超出 CAP-01—CAP-06，先按 implementation-discovery 记录；不得直接创建未登记能力。
7. 不执行推送、tag、release、部署、外部消息、数据迁移或生产命令；遇到 host/non-offline action 必须返回外部动作 Gate。

## 波次与依赖

```text
W0 baseline/owner Gate
  ↓
W1 public model + error/result contracts
  ↓
W2 loader/render/condition + graph validation
  ↓
W3 process/file/handler engines + cancellation/output
  ↓
W4 runner integration + standalone consumer tests
  ↓
W5 Kite legacy adapter and compatibility tests
  ↓
W6 second consumer + Go 1.23 matrix + documentation
```

W1—W3 可在目标 module 内顺序推进，但每个波次完成后先做 focused validation。W4 之前不能改 Kite 调用方；W5 之前不能声称兼容旧配置；W6 之前不能声称跨应用可用。

## 任务

### T01 目标仓库与基线确认

- 文件: `D:/work/inhere/my-tools-dev/gookit2/kscript/.git`（目标 Git root，当前尚不存在）；现有 `inhere-tools/kite-go` 仅读 baseline 文件。
- 动作: owner 确认 `gookit2/kscript` 目录、Git root、Go 1.23 CI 方案和第二应用；创建目标 module 前记录 branch/HEAD/status；不要触碰 Kite 的 `.codebase-memory/`、`docs/research/`。
- 验证: `git -C <target> status --short`、`go env GOVERSION`、目标模块 `go list ./...`（创建后）；保存 baseline/progress 记录。
- 完成标准: 目标 root 和 owner 有明确证据；无未分类 dirty 路径；D02—D10 没有新的核心歧义。
- 依赖: 设计批准、计划批准；本任务之前不创建文件。

### T02 创建独立 module 和公共模型

- 文件: `<target>/go.mod`、`<target>/kscript.go`、`<target>/definition.go`、`<target>/request.go`、`<target>/result.go`、`<target>/engine.go`、`<target>/errors.go`。
- 动作: 设置 module `github.com/gookit/kscript`、`go 1.23`；实现 Definition/Task/Step/ExecSpec/ShellSpec/FileSpec/HostSpec、Request/IO、Result/RunError、Engine/Handler；复制输入并冻结 Runner Definition；不导入 Kite/internal、gcli、cliui、全局日志。
- 验证: `go test ./...`；外部临时 module `go 1.23` 使用 replace 只做本地编译；`go list -deps` 检查无 `github.com/inhere/kite-go`。
- 完成标准: 公共类型有 GoDoc；无效空 Task、负 timeout、未知动作和可变输入返回稳定错误；基础 exec Definition 可构造。
- 依赖: T01。

### T03 实现格式 loader、来源和显式发现

- 文件: `<target>/formats/doc.go`、`<target>/formats/yaml.go`、`<target>/formats/json.go`、`<target>/formats/toml.go`、`<target>/formats/discover.go`、`<target>/internal/validate/*`。
- 动作: 实现 schema version 1 解码到统一 Definition；记录 Source name/BaseDir；重复 key/未知字段/来源冲突/路径越界失败；实现 nearest/ancestors、StartDir/StopDir/MaxDepth 显式选项；不默认读磁盘或父目录。
- 验证: 同一 fixture 的三种格式模型相等；非法 fixture 均在 Engine 前失败；发现顺序、边界和 Windows 路径测试；`go test -race ./...`。
- 完成标准: Loader 无命令和 Handler 副作用，Definition 只在完整校验后发布；来源可定位到文件和行/字段。
- 依赖: T02。

### T04 实现变量、表达式和条件判断

- 文件: `<target>/internal/render/*`、`<target>/internal/condition/*`、`<target>/definition.go`、`<target>/inspect.go`。
- 动作: 实现命名空间插值、Vars/Env/Args/HostData/运行元数据只读视图；接入 expr；Task `if` 在 deps 前求值，Step `if` 在 dynamic_vars 后、Engine 前求值；bool/false/Deferred/编译错误语义按设计；动态变量不得在 Inspect 执行。
- 验证: `go test ./internal/render ./internal/condition`；true/false/空/非 bool/未知变量/动态条件 fixture；Inspect 副作用计数为 0；`go test -race ./...`。
- 完成标准: 条件不 panic、不打印、不固定返回 true；任务级和 Step 级 skipped 状态及原因在 Result/Plan 中可观察。
- 依赖: T02、T03。

### T05 实现顺序图调度和运行状态

- 文件: `<target>/internal/graph/*`、`<target>/runner.go`、`<target>/call_state.go`。
- 动作: 预检 deps/task call 全图、报告完整环路径、限制最大深度/展开数；按出现次数串行执行，不隐式去重；每个 CallID 独立 Vars/Env/Dir/deadline；Task/Step 条件和 platform skip 语义；ignore_error 只容忍动作错误。
- 验证: 缺失/环/膨胀/菱形/重复调用、变量污染、依赖失败阻断和 skipped fixture；顺序事件断言；`go test -race ./...`。
- 完成标准: 同一 Runner 并发 Run 无数据串扰；所有定义错误在首个动作前返回。
- 依赖: T02、T03、T04。

### T06 实现进程、Shell、file、host Engine

- 文件: `<target>/internal/process/*`、`<target>/engine_process.go`、`<target>/engine_file.go`、`<target>/handler.go`。
- 动作: exec 使用 argv；Shell 结构化选择 sh/bash/zsh/cmd/pwsh；file 使用 Program/PrefixArgs；Handler 显式注册且不可在 New 后改变；实现 stdin/stdout/stderr、CaptureLimit、错误包装、dry-run；不调用 os.Chdir/os.Setenv。
- 验证: Linux/macOS Shell（可用时）和 Windows cmd/pwsh 由 CI/host matrix 执行；参数空格/引号；解释器缺失；Handler context；dry-run zero side-effect；大输出和 writer error。
- 完成标准: Process Engine 返回启动/退出/取消/超时/I/O 分类；shell 源码和 exec argv 不互相隐式转换。
- 依赖: T02、T05。

### T07 实现取消、超时和进程清理

- 文件: `<target>/internal/process/*_unix.go`、`<target>/internal/process/*_windows.go`、`<target>/result.go`、`<target>/runner.go`。
- 动作: Root/Task/Step deadline 取最早值；取消停止新步骤；POSIX 进程组、Windows Job Object 或已验证等价方案清理所拥有子进程；Handler 只遵守协作取消；保留 context 原因。
- 验证: 子进程树 fixture、取消/超时/宽限期、无逃逸回收、Go 1.23 Windows/Linux matrix；`go test -race ./...`。
- 完成标准: 不静默只杀父进程；能力不可用时启动前明确失败；Result 状态为 Canceled 或 TimedOut。
- 依赖: T06。

### T08 集成 Runner、Inspect、Result 和示例

- 文件: `<target>/runner.go`、`<target>/inspect.go`、`<target>/examples/basic/main.go`、`<target>/examples/config/main.go`、`<target>/README.md`、`<target>/README.zh-CN.md`。
- 动作: 暴露 New/Lookup/List/Inspect/Run；补齐 GoDoc、配置 schema、错误表、最小示例和迁移说明；记录 stdout/stderr、步骤状态和事件；不把示例中的 log.Fatal 变成库行为。
- 验证: `go test ./...`、`go vet ./...`、`go test -race ./...`；从外部临时 module 以 `go 1.23` 引入；README 示例命令实际运行。
- 完成标准: 外部应用只导入主包即可执行基础任务；Inspect 不执行副作用；API 与设计保持一致。
- 依赖: T02、T04、T05、T06、T07。

### T09 Kite legacy converter 和调用方迁移

- 文件: `kite-go/internal/biz/cmdbiz/runany.go`、`kite-go/internal/boot/service.go`、`kite-go/internal/cli/*`（实际变更前重新确认）、`<target>/formats/legacy_kite.go`、`kite-go/go.mod`。
- 动作: 在 Kite 侧将旧 string/list/map、`__settings`、`@task:`、`$@/$*`、`@sh:/@exec:`、ScriptDirs 和 `AppendVarsFn` 转为新 Definition/Request；保留 alias/ext/plugin/system 外层优先级；分阶段替换 import，禁止把 Kite internal 迁入新库。
- 验证: 旧 fixture 前后对照；Kite `go test ./...`；RunAny 对根 ErrNotFound 才系统命令兜底，执行错误不吞；Kite 与新库分别 `go list -deps` 检查边界。
- 完成标准: 合法旧任务行为达到设计兼容清单；修复差异逐项记录；任何新增生效字段都有迁移说明。
- 依赖: T03、T04、T05、T06、T08；若发现行为/范围变化，先回 Design Gate。

### T10 第二应用 consumer 验证和 Go 1.23+ 矩阵

- 文件: `<consumer-project>/go.mod`、`<consumer-project>/...`（由 consumer owner 指定）；`<target>/.github/workflows/*` 或现有 CI 配置。
- 动作: 使用真实第二应用接入 Task/Step if、Handler、取消和输出；补 Go 1.23 与当前稳定版本矩阵；不把临时示例替代真实接入。
- 验证: consumer 项目自身构建、单元/集成测试；`go test -race ./...`；矩阵日志保存版本、module 版本和实际 Result；`go list -m all` 检查依赖。
- 完成标准: consumer 无 Kite 初始化即可运行；Go 1.23 通过；不使用高版本专属 API；问题按 Operational Discovery/Corrective/Semantic Amendment 分类。
- 依赖: T08；第二应用路径必须在 T01 Gate 确认。

### T11 文档、版本和发布准备审查

- 文件: `<target>/README.md`、`<target>/README.zh-CN.md`、`<target>/CHANGELOG.md`、`<target>/LICENSE`、`<target>/docs/*`、`kite-go/docs/*`。
- 动作: 固定 API/配置/迁移/安全/平台限制；保留 MIT 版权；记录 v0.x 兼容政策、Go 1.23+、条件判断和后续能力边界；运行 license/dependency 检查；形成 release candidate 但不发布。
- 验证: Markdown 链接检查、`go vet ./...`、`go test -race ./...`、外部 consumer 再编译；审查 dirty/status 和只 stage 所有者文件。
- 完成标准: 文档与实际 API/测试一致；无未分类 secret/依赖或平台限制；发布候选等待单独批准。
- 依赖: T09、T10。

## 回滚与恢复

- 每个波次完成后建立本地 atomic commit，并记录 commit、scope、验证命令和状态；不 stage `.codebase-memory/`、`docs/research/` 或其他未归属路径。
- 新 module 回滚优先删除未发布的适配提交并恢复 Kite 原 import；不自动重写或删除用户任务文件。
- Kite 切换期间保留可配置旧实现开关或独立回退 commit；回退只恢复代码/依赖，不撤销已运行脚本产生的文件、服务调用或外部副作用。
- 任意发现 `Semantic Amendment`、`Ownership Conflict`、目标路径冲突、API 破坏性变化或无法证明子进程清理时，在下一次 mutation 前暂停，更新 progress 并返回 Design/Plan review Gate。
- release、push、tag、部署、迁移和外部动作不属于自动回滚；没有对应授权时保持本地候选状态。

## 人工 Gate

1. **设计批准 Gate**：批准 Draft 0.2、Go 1.23+、Task/Step 条件判断和 v0.1 范围；未批准不得进入 T01 之后的 mutation。
2. **计划批准 Gate**：批准本计划的目标 module 路径 `gookit2/kscript`、T01—T11 任务和 `host_or_non_offline_action=NOT_APPLICABLE`。
3. **目标仓库/consumer Gate**：在 T01 明确新 module Git root、owner、第二应用路径和 Go 1.23 runner；路径或 owner 不明即停。
4. **实施 Gate**：计划批准后仍需用户单独发出当前执行请求；该请求到达前只可审查或修订文档。
5. **外部动作 Gate**：若后续要创建远端仓库、推送、tag、release、部署或运行非离线动作，逐项列出目标并另行取得批准。

## 可追溯性

| 设计目标/验收 | 任务 | 主要验证 |
|---|---|---|
| 独立 module 与 Go 1.23+ | T01, T02, T08, T10 | 外部 module 编译、Go 1.23 matrix、无 Kite 依赖 |
| 类型化 API 和不可变运行状态 | T02, T05, T08 | API 编译、输入隔离、race、并发 Run |
| YAML/JSON/TOML 与显式发现 | T03 | 等价 fixture、非法配置、来源和发现顺序 |
| 任务/Step 条件判断 | T04, T05, T08 | expr bool、Task/Step 时机、Skipped/Deferred、零副作用 Inspect |
| deps/call 和错误边界 | T05, T09 | 循环/膨胀/重复调用、ErrNotFound 与执行错误分类 |
| exec/Shell/file/host | T06, T09, T10 | argv/Shell/file matrix、Handler、Kite adapter |
| 取消、超时、输出和清理 | T06, T07, T08, T10 | 进程树、deadline、CaptureLimit、writer error |
| Kite 兼容和第二应用结果 | T09, T10 | 旧 fixture、Kite tests、真实 consumer Result |
| 文档、生命周期和回滚 | T01, T09, T11 | status/stage、license、迁移文档、candidate review |

## 完成 Gate 与剩余工作

计划完成需要：T01—T11 各任务有 owner、实际路径、验证命令和 progress；目标 module 与第二应用已确认；所有 focused tests、race、Go 1.23 matrix、Kite fixture 和 external consumer 证据通过；无开放 Semantic Amendment/Ownership Conflict；只 stage 归属文件并完成原子本地提交。

完成 Gate 不包括 push、tag、release、部署、数据迁移或外部消息。发布前还需独立版本/发布审查和外部动作批准。并行、重试、缓存、Taskfile/justfile、mvdan/sh 和脚本 VM 继续作为后续设计，不得在本计划执行中以“顺手完善”扩大范围。

本计划生成后停止在计划批准 Gate；计划本身不授权实施。
