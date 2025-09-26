# Kite-Go 项目开发指南

## 项目概述

`kite-go` 是一个功能强大的个人开发者工具集合，基于 Go 语言开发的命令行应用程序。项目旨在为开发者提供一套完整的开发、部署、维护工具。

### 核心特性

- **Git 工具集**: 提供 Git、GitLab、GitHub 常用操作的封装
- **文本处理**: 强大的字符串、JSON、YAML 等格式处理工具
- **文件系统**: 文件查看、渲染、模板处理等功能
- **HTTP 工具**: HTTP 服务器、API 测试、模板请求等
- **系统工具**: 环境变量管理、剪贴板操作、可执行文件搜索等
- **脚本执行**: 支持多种脚本类型和任务管理
- **快速跳转**: 目录历史记录和快速导航
- **插件扩展**: 支持动态插件和扩展机制

## 项目架构

### 目录结构

```
kite-go/
├── cmd/                    # 可执行程序入口
│   ├── kite/              # 主程序入口
│   ├── htu/               # HTTP 工具
│   ├── ktenv/             # 环境管理工具
│   └── pacgo/             # PAC 工具
├── internal/              # 内部包，不对外暴露
│   ├── app/               # 应用核心
│   ├── appconst/          # 应用常量
│   ├── apputil/           # 应用工具函数
│   ├── biz/               # 业务逻辑
│   ├── bootstrap/         # 启动引导
│   ├── cli/               # CLI 命令实现
│   │   ├── aicmd/         # AI 相关命令
│   │   ├── appcmd/        # 应用管理命令
│   │   ├── devcmd/        # 开发工具命令
│   │   ├── fscmd/         # 文件系统命令
│   │   ├── gitcmd/        # Git 相关命令
│   │   ├── httpcmd/       # HTTP 工具命令
│   │   ├── syscmd/        # 系统命令
│   │   ├── textcmd/       # 文本处理命令
│   │   └── toolcmd/       # 工具命令
│   ├── initlog/           # 日志初始化
│   └── web/               # Web 相关
├── pkg/                   # 可重用包
│   ├── aitool/            # AI 工具
│   ├── cmdutil/           # 命令工具
│   ├── common/            # 通用功能
│   ├── gitx/              # Git 扩展
│   ├── httptpl/           # HTTP 模板
│   ├── kiteext/           # Kite 扩展
│   ├── kscript/           # 脚本引擎
│   ├── quickjump/         # 快速跳转
│   └── simpleai/          # 简单 AI 集成
├── config/                # 配置文件
├── data/                  # 数据文件
├── static/                # 静态资源
└── test/                  # 测试文件
```

### 核心组件

#### 1. 应用核心 (`internal/app/`)

- **`app.go`**: 应用主入口和生命周期管理
- **`appconf.go`**: 配置管理
- **`boot.go`**: 启动引导逻辑
- **`service.go`**: 服务管理

#### 2. 命令行系统 (`internal/cli/`)

采用分组式命令结构：

- **Git 工具组**: `git`, `gitlab`, `github` 命令
- **文件系统工具组**: `fs` 命令
- **文本处理工具组**: `text`, `json` 命令
- **HTTP 工具组**: `http` 命令
- **系统工具组**: `sys` 命令
- **开发工具组**: `dev` 命令

#### 3. 扩展系统 (`pkg/kiteext/`)

- **插件管理**: 动态加载和执行插件
- **脚本引擎**: 支持多种脚本格式
- **路径映射**: 路径别名和解析
- **变量映射**: 动态变量替换

#### 4. 启动系统 (`internal/bootstrap/`)

- **配置加载**: 多层级配置文件加载
- **服务初始化**: 各种服务组件初始化
- **CLI 构建**: 命令行应用构建

## 开发环境搭建

### 基础要求

- **Go 版本**: 1.24+ (项目使用 Go 1.24.6 工具链)
- **操作系统**: 支持 Windows、Linux、macOS
- **Git**: 版本控制

### 开发依赖

项目使用 Go Modules 管理依赖，核心依赖包括：

#### CLI 框架
```go
github.com/gookit/gcli/v3 v3.2.3        // CLI 应用框架
```

#### 配置管理
```go
github.com/gookit/config/v2 v2.2.5      // 配置文件管理
github.com/gookit/ini/v2 v2.3.1         // INI 配置支持
github.com/goccy/go-yaml v1.18.0        // YAML 支持
```

#### 文本处理
```go
github.com/alecthomas/chroma v0.10.0     // 语法高亮
github.com/charmbracelet/glamour v0.10.0 // Markdown 渲染
github.com/gomarkdown/markdown v0.0.0-20250810172220-2e2c11897d1a // Markdown 处理
```

#### 模板引擎
```go
github.com/CloudyKit/jet/v6 v6.3.1      // Jet 模板引擎
github.com/gookit/easytpl v1.1.0        // 简单模板
```

#### HTTP 和 Web
```go
github.com/gookit/rux v1.4.0            // HTTP 路由器
github.com/gookit/greq v0.4.0           // HTTP 客户端
```

#### 脚本执行
```go
github.com/traefik/yaegi v0.16.1        // Go 脚本解释器
github.com/expr-lang/expr v1.17.6       // 表达式求值
```

#### AI 集成
```go
github.com/sashabaranov/go-openai v1.41.1 // OpenAI API 客户端
```

#### 工具库
```go
github.com/gookit/goutil v0.7.1         // 通用工具库
github.com/gookit/color v1.6.0          // 颜色输出
github.com/gookit/slog v0.5.8           // 结构化日志
github.com/gookit/gitw v0.3.5           // Git 包装器
```

### 环境设置

1. **克隆项目**
```bash
git clone https://github.com/inhere/kite-go.git
cd kite-go
```

2. **安装依赖**
```bash
go mod download
```

3. **设置开发环境变量**
```bash
# 调试模式
export KITE_VERBOSE=debug

# Windows PowerShell
$env:KITE_VERBOSE='debug'
```

4. **构建和运行**
```bash
# 开发模式运行
go run ./cmd/kite

# 构建
make build

# 安装到 GOPATH/bin
make install
```

## 开发流程

### 1. 项目初始化流程

```
main()
  ↓
bootstrap.MustRun(app.App())
  ↓
app.Boot() - 启动应用
  ↓
runBootloaders() - 执行启动加载器
  ↓
cli.Boot() - 加载命令
  ↓
app.Run() - 运行应用
```

### 2. 添加新命令

1. **创建命令文件**
在 `internal/cli/` 相应的命令组目录下创建新的命令文件：

```go
// internal/cli/newcmd/example.go
package newcmd

import (
    "github.com/gookit/gcli/v3"
)

var ExampleCmd = &gcli.Command{
    Name: "example",
    Desc: "示例命令描述",
    Config: func(c *gcli.Command) {
        // 配置选项和参数
        c.AddArg("name", "名称参数", true)
        c.BoolOpt2(&opts.verbose, "verbose,v", "详细输出")
    },
    Func: func(c *gcli.Command, args []string) error {
        // 命令执行逻辑
        return nil
    },
}
```

2. **注册命令**
在 `internal/cli/boot.go` 的 `addCommands()` 函数中注册：

```go
func addCommands(cli *gcli.App) {
    cli.Add(
        // ... 其他命令
        newcmd.ExampleCmd,
    )
}
```

### 3. 添加新的扩展包

1. **创建包目录**
在 `pkg/` 下创建新的包目录

2. **实现核心接口**
根据需要实现相应的接口，如：
- `BootLoader` - 启动加载器
- `Extension` - 扩展接口

3. **注册扩展**
在 `internal/bootstrap/service.go` 中注册新的服务

### 4. 配置管理

配置文件位于 `config/` 目录：
- `config.yml` - 主配置文件
- `config.windows.yml` - Windows 特定配置
- `config.darwin.yml` - macOS 特定配置

配置加载顺序：
1. 默认配置
2. 基础配置文件
3. 平台特定配置
4. 用户自定义配置
5. 环境变量覆盖

### 5. 测试

运行测试：
```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./internal/cli/...

# 运行单个测试文件
go test ./test/unittest/cli/
```

## 构建和部署

### 本地构建

```bash
# 构建当前平台
make build

# 构建所有平台
make build-all

# 构建特定平台
make linux      # Linux AMD64
make win        # Windows AMD64
make darwin     # macOS AMD64
```

### 发布

使用 GitHub Actions 自动发布：
- 推送 tag 触发自动构建
- 支持多平台交叉编译
- 自动创建 GitHub Release

### 安装脚本

提供快速安装脚本：
```bash
curl https://raw.githubusercontent.com/inhere/kite-go/main/cmd/install.sh | bash
```

## 最佳实践

### 1. 代码组织

- **`internal/`**: 应用私有代码，不对外暴露
- **`pkg/`**: 可重用的公共包
- **`cmd/`**: 可执行程序入口
- 每个命令组有独立的包

### 2. 错误处理

使用 `github.com/gookit/goutil/errorx` 进行错误包装：

```go
if err != nil {
    return errorx.Wrapf(err, "操作失败: %s", operation)
}
```

### 3. 日志记录

使用结构化日志：

```go
app.Log().WithValue("key", value).Info("操作完成")
```

### 4. 配置访问

通过应用实例访问配置：

```go
config := app.Cfg()
value := config.String("key.subkey")
```

### 5. 命令设计

- 遵循 Unix 哲学：做好一件事
- 提供清晰的帮助信息
- 支持管道操作
- 合理的参数和选项设计

## 贡献指南

### 开发规范

1. **代码风格**: 遵循 Go 官方代码规范
2. **提交信息**: 使用约定式提交 (Conventional Commits)
3. **测试覆盖**: 新功能必须包含测试
4. **文档更新**: 重要功能需要更新文档

### 提交流程

1. Fork 项目
2. 创建功能分支
3. 编写代码和测试
4. 提交 Pull Request

### 调试技巧

1. **启用调试模式**:
```bash
export KITE_VERBOSE=debug
```

2. **查看应用信息**:
```bash
kite app info
```

3. **查看配置**:
```bash
kite app config -a
```

4. **查看路径信息**:
```bash
kite app path -a
```

## 相关资源

- **项目主页**: https://github.com/inhere/kite-go
- **文档**: 项目 `docs/` 目录
- **问题反馈**: GitHub Issues
- **Gookit 工具库**: https://github.com/gookit

---

此文档为开发者提供了完整的项目理解和开发指导，帮助快速上手 kite-go 项目的开发工作。