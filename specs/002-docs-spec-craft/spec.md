# Feature Specification: Kite XEnv 命令行工具

**Feature Branch**: `002-docs-spec-craft`  
**Created**: 2025年10月16日  
**Status**: Draft  
**Input**: User description: "请从 @docs/spec-craft/kite-xenv-spec-craft.md 查看新功能描述并实现"

## Clarifications

### Session 2025-10-16

- Q: 用户角色应如何区分？ → A: 仅单一用户类型
- Q: 应记录何种详细程度的日志？ → A: 记录基本操作日志，不含敏感信息
- Q: 配置导出的大小应如何限制？ → A: 限制为10MB
- Q: 错误信息应显示何种详细程度？ → A: 详细错误信息
- Q: 卸载工具时如何处理相关配置？ → A: 用户自行决定

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 安装和管理开发工具链 (Priority: P1)

开发者需要在系统上安装、更新和卸载不同版本的开发工具（如Go、Node.js、Java等），以便在不同项目中使用合适的工具版本。

**Why this priority**: 这是xenv工具的核心功能，没有这个功能，其他功能都没有意义。

**Independent Test**: 开发者可以使用`kite xenv tools install <name:version>`命令安装特定版本的工具，并使用`kite xenv tools list`验证安装是否成功。

**Acceptance Scenarios**:

1. **Given** 开发者已经安装了kite，**When** 执行`kite xenv tools install go@1.21`，**Then** 系统下载并安装Go 1.21版本，并在工具链列表中显示
2. **Given** 某个工具版本已经安装，**When** 执行`kite xenv tools uninstall go@1.21`，**Then** 系统移除该工具版本并更新工具列表
3. **Given** 工具链已经安装，**When** 执行`kite xenv tools update go@latest`，**Then** 系统更新到最新可用版本

---

### User Story 2 - 切换和激活工具链版本 (Priority: P1)

开发者需要在不同项目中切换和激活不同版本的开发工具，以匹配项目需求。

**Why this priority**: 这是开发人员日常工作中最常用的功能，直接关系到开发效率。

**Independent Test**: 开发者可以使用`kite xenv use go@1.21`命令激活特定版本的工具，并通过运行命令验证是否生效。

**Acceptance Scenarios**:

1. **Given** 多个版本的工具已安装，**When** 执行`kite xenv use go@1.21`，**Then** 当前会话中Go命令指向1.21版本
2. **Given** 某个工具版本已激活，**When** 使用`kite xenv use go@1.20 -g`进行全局激活，**Then** 所有新会话都使用Go 1.20版本
3. **Given** 工具版本已全局激活，**When** 使用`kite xenv unuse go@1.20 -g`取消激活，**Then** 全局设置中不再包含该工具版本

---

### User Story 3 - 管理环境变量和PATH路径 (Priority: P2)

开发者需要设置和管理环境变量和PATH路径，以支持不同的开发需求。

**Why this priority**: 这是开发环境管理的重要组成部分，许多工具需要特定的环境变量或PATH设置才能正常工作。

**Independent Test**: 开发者可以使用`kite xenv env --set`命令设置环境变量，并使用`kite xenv list --env`验证设置是否生效。

**Acceptance Scenarios**:

1. **Given** 用户需要设置环境变量，**When** 执行`kite xenv env --set NODE_ENV production`，**Then** 当前会话中NODE_ENV变量被设置为production
2. **Given** 用户需要添加PATH路径，**When** 执行`kite xenv path --add ~/.custom-tools/bin`，**Then** PATH环境变量中包含新路径
3. **Given** 某个环境变量已设置，**When** 执行`kite xenv env --unset NODE_ENV`，**Then** 该环境变量被移除

---

### User Story 4 - 导入导出配置 (Priority: P3)

开发者需要在不同机器间同步开发环境配置，以便快速设置新的开发环境。

**Why this priority**: 这个功能提高了开发人员在多个设备间的迁移效率，特别是在团队协作中很有价值。

**Independent Test**: 开发者可以使用`kite xenv config --export`导出配置，然后在另一台机器上使用`kite xenv config --import`导入。

**Acceptance Scenarios**:

1. **Given** 用户已配置好开发环境，**When** 执行`kite xenv config --export zip`，**Then** 生成包含所有配置的ZIP文件
2. **Given** 有导出的配置文件，**When** 在新机器上执行`kite xenv config --import config.zip`，**Then** 系统恢复所有配置的工具链和环境设置

---

### User Story 5 - Shell集成和实时生效 (Priority: P2)

开发者需要在切换工具链后立即生效，无需重启shell或执行额外命令。

**Why this priority**: 实时生效是现代环境管理工具的基本要求，提高了开发效率。

**Independent Test**: 开发者配置shell hook `kite xenv shell --type bash` 后，切换工具版本应立即在当前shell中生效。

**Acceptance Scenarios**:

1. **Given** 用户已配置shell hook，**When** 执行`kite xenv use go@1.21`，**Then** 当前shell中立即使用Go 1.21版本
2. **Given** 用户在项目目录下有配置文件，**When** 进入目录时，**Then** 自动应用该目录的环境设置

---

### Edge Cases

- 如果网络不可用时尝试下载工具会如何处理？
- 如何处理磁盘空间不足时安装工具的情况？
- 当多个用户共享同一系统时，如何隔离各自的环境配置？
- 当用户尝试安装不存在的工具版本时，系统如何响应？
- 卸载工具时系统应如何处理相关配置？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须允许用户安装、卸载和更新多种开发工具的不同版本
- **FR-002**: 系统必须支持从外部源（如GitHub）下载和安装工具
- **FR-003**: 用户必须能够通过`kite xenv use`命令切换和激活不同的工具链版本
- **FR-004**: 系统必须管理用户的环境变量和PATH路径
- **FR-005**: 系统必须允许用户导入和导出配置以实现跨机器同步
- **FR-006**: 系统必须提供shell hooks以实现实时环境切换
- **FR-007**: 系统必须支持PowerShell、Bash和Zsh等主流shell
- **FR-008**: 系统必须兼容执行目录级配置文件（如`.envrc`, `.xenv.toml`）
- **FR-009**: 系统必须支持全局和会话级别的环境设置
- **FR-010**: 系统必须在Windows、macOS和Linux上运行并保持一致的行为
- **FR-011**: 系统必须记录基本操作日志，不含敏感信息
- **FR-012**: 系统在错误情况下必须提供详细错误信息
- **FR-013**: 系统在卸载工具时必须允许用户自行决定是否保留相关配置

### Key Entities *(include if feature involves data)*

- **Tool Chain**: 代表特定版本的开发工具（如Go、Node.js等），包含版本号、安装路径、别名等属性
- **Environment Variable**: 代表系统中的环境变量，具有名称、值、作用域（全局/会话）属性
- **Configuration**: 代表用户的配置信息，包含工具管理设置、路径配置、环境激活状态等数据
- **Path Entry**: 代表添加到PATH环境变量中的路径条目，具有路径值和优先级属性
- **User**: 代表使用xenv工具的单一用户类型，具有对自身配置和环境的完全控制权

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 开发者能够在5分钟内完成特定版本开发工具的安装和激活
- **SC-002**: 系统能够在所有支持的平台上正确运行，实现跨平台配置同步
- **SC-003**: 工具链切换命令应在1秒内完成激活，对用户透明生效
- **SC-004**: 配置导出和导入功能应支持所有环境设置（工具版本、环境变量、PATH），且单次导出大小不超过10MB
- **SC-005**: 95%的用户在使用shell hooks后，无需重启shell即可立即使用新激活的工具
