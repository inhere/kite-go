# Project Context

## Purpose
Kite 是一个个人开发者工具命令行应用程序，旨在提供一套统一的开发工具集，包括 Git 操作、文件处理、HTTP 服务、系统信息查看等功能。它简化了日常开发工作流程，提供了快速命令别名和脚本执行能力。

## Tech Stack
- **Go 1.24** - 主要编程语言
- **gookit/gcli/v3** - CLI 框架
- **gookit/config/v2** - 配置管理
- **gookit/rux** - HTTP 路由
- **gookit/gitw** - Git 操作封装
- **gookit/goutil** - 通用工具库
- **gomarkdown/markdown** - Markdown 处理
- **chroma** - 语法高亮
- **charmbracelet/glamour** - 终端渲染

## Project Conventions

### Code Style
- 使用 Go 标准格式化工具 `go fmt`
- 使用 `goimports` 管理导入
- 遵循 Go 官方代码规范
- 使用 `make lint` 进行代码检查
- 变量命名采用驼峰命名法
- 常量使用大写字母和下划线

### Architecture Patterns
- **模块化设计**: 按功能模块组织代码（pkg/ 目录）
- **命令模式**: 每个功能作为独立的命令实现
- **配置驱动**: 通过 YAML 配置文件管理应用行为
- **插件系统**: 支持扩展插件和自定义脚本
- **分层架构**: internal/ 目录包含内部实现，pkg/ 目录包含可重用组件

### Testing Strategy
- 单元测试位于 `test/unittest/` 目录
- 集成测试位于 `test/integration/` 目录
- 使用 Go 标准测试框架
- CI/CD 通过 GitHub Actions 自动运行测试
- 测试覆盖率要求：核心功能 >= 80%

### Git Workflow
- **主分支**: `main` - 稳定版本
- **开发模式**: 功能开发在 `main` 分支进行
- **提交规范**: 使用语义化提交信息（feat:, fix:, docs:, etc.）
- **版本管理**: 基于 Git 标签进行版本控制
- **自动化**: 通过 GitHub Actions 进行 CI/CD

## Domain Context
Kite 是一个面向开发者的多用途命令行工具，主要服务于：
- **日常开发工作流**: Git 操作、文件处理、目录跳转
- **API 测试**: HTTP 请求发送、服务模拟
- **系统管理**: 环境信息查看、命令执行
- **脚本自动化**: 自定义脚本执行和任务管理
- **文档处理**: Markdown 渲染、模板处理

## Important Constraints
- **跨平台兼容性**: 支持 Windows、Linux、macOS
- **Go 1.24+**: 最低 Go 版本要求
- **轻量级设计**: 单一二进制文件，无外部依赖
- **配置优先**: 通过配置文件而非代码修改定制行为
- **向后兼容**: API 变更需要保持向后兼容性

## External Dependencies
- **Git**: 版本控制系统（必需）
- **Shell 环境**: 用于命令执行和脚本运行
- **网络连接**: 用于 HTTP 请求和远程操作
- **文件系统**: 用于文件操作和配置管理
- **终端环境**: 用于交互式操作和显示
