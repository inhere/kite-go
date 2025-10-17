---
description: "Task list for Kite XEnv implementation"
---

# Tasks: Kite XEnv 命令行工具

**Input**: Design documents from `/specs/002-docs-spec-craft/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `cmd/`, `pkg/`, `tests/` at repository root
- **CLI Application**: `cmd/kite/main.go`, `pkg/xenv/` per plan.md structure
- Paths shown below follow the project structure from plan.md

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create pkg/xenv directory structure per implementation plan
- [x] T002 Create pkg/xenv/config directory for configuration management
- [x] T003 [P] Create pkg/xenv/tools directory for tool chain management
- [x] T004 [P] Create pkg/xenv/env directory for environment variable management
- [x] T005 [P] Create pkg/xenv/shell directory for shell integration hooks
- [x] T006 [P] Create pkg/xenv/models directory for data models
- [x] T007 Create tests/integration/xenv directory for integration tests

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T008 Create base ToolChain model in pkg/xenv/models/toolchain.go based on data-model.md
- [x] T009 Create base EnvironmentVariable model in pkg/xenv/models/env_var.go based on data-model.md
- [x] T010 [P] Create base PathEntry model in pkg/xenv/models/path_entry.go based on data-model.md
- [x] T011 [P] Create base Configuration model in pkg/xenv/models/config.go based on data-model.md
- [x] T012 [P] Create base ActivityState model in pkg/xenv/models/activity_state.go based on data-model.md
- [x] T013 [P] Create base User model in pkg/xenv/models/user.go based on data-model.md
- [x] T014 Implement configuration management in pkg/xenv/config/config.go
- [x] T015 Setup CLI command structure in cmd/kite/main.go using gcli framework
- [x] T016 [P] Implement file system utility functions in internal/util/fs.go
- [x] T017 Implement cross-platform path utilities in internal/util/path.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - 安装和管理开发工具链 (Priority: P1) 🎯 MVP

**Goal**: 实现安装、卸载和更新开发工具的功能

**Independent Test**: 开发者可以使用`kite xenv tools install <name:version>`命令安装特定版本的工具，并使用`kite xenv tools list`验证安装是否成功

### Tests for User Story 1 (OPTIONAL - only if tests requested) ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T018 [P] [US1] Contract test for `kite xenv tools install <name:version>` in tests/integration/xenv/test_tools_install.go
- [ ] T019 [P] [US1] Contract test for `kite xenv tools uninstall <name:version>` in tests/integration/xenv/test_tools_uninstall.go
- [ ] T020 [P] [US1] Contract test for `kite xenv tools update <name:version>` in tests/integration/xenv/test_tools_update.go

### Implementation for User Story 1

- [x] T021 [P] [US1] Implement ToolChain service in pkg/xenv/tools/service.go
- [x] T022 [US1] Implement tool installation logic in pkg/xenv/tools/installer.go
- [x] T023 [US1] Implement tool download functionality in pkg/xenv/tools/downloader.go
- [x] T024 [US1] Implement tool uninstallation logic in pkg/xenv/tools/uninstaller.go
- [x] T025 [US1] Implement tool listing in pkg/xenv/tools/list.go
- [x] T026 [US1] Create xenv tools command in cmd/kite/xenv_tools_cmd.go
- [x] T027 [US1] Add tools install subcommand to tools command
- [x] T028 [US1] Add tools uninstall subcommand to tools command
- [x] T029 [US1] Add tools update subcommand to tools command
- [x] T030 [US1] Add tools list subcommand to tools command
- [x] T031 [US1] Add tools show subcommand to tools command

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - 切换和激活工具链版本 (Priority: P1)

**Goal**: 实现切换和激活不同版本的开发工具，以匹配项目需求

**Independent Test**: 开发者可以使用`kite xenv use go@1.21`命令激活特定版本的工具，并通过运行命令验证是否生效

### Tests for User Story 2 (OPTIONAL - only if tests requested) ⚠️

- [ ] T032 [P] [US2] Contract test for `kite xenv use <name:version>` in tests/integration/xenv/test_use_cmd.go
- [ ] T033 [P] [US2] Contract test for `kite xenv unuse <name:version>` in tests/integration/xenv/test_unuse_cmd.go

### Implementation for User Story 2

- [x] T034 [P] [US2] Implement tool activation logic in pkg/xenv/tools/activator.go
- [x] T035 [US2] Implement tool deactivation logic in pkg/xenv/tools/deactivator.go
- [x] T036 [US2] Update ActivityState model to track active tools
- [x] T037 [US2] Create xenv use command in cmd/kite/xenv_use_cmd.go
- [x] T038 [US2] Create xenv unuse command in cmd/kite/xenv_unuse_cmd.go
- [x] T039 [US2] Implement global vs session scope handling for tool activation

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - 管理环境变量和PATH路径 (Priority: P2)

**Goal**: 实现设置和管理环境变量和PATH路径，以支持不同的开发需求

**Independent Test**: 开发者可以使用`kite xenv env --set`命令设置环境变量，并使用`kite xenv list --env`验证设置是否生效

### Tests for User Story 3 (OPTIONAL - only if tests requested) ⚠️

- [ ] T040 [P] [US3] Contract test for `kite xenv env --set <name> <value>` in tests/integration/xenv/test_env_set.go
- [ ] T041 [P] [US3] Contract test for `kite xenv env --unset <name>` in tests/integration/xenv/test_env_unset.go
- [ ] T042 [P] [US3] Contract test for `kite xenv path --add <path>` in tests/integration/xenv/test_path_add.go
- [ ] T043 [P] [US3] Contract test for `kite xenv path --rm <path>` in tests/integration/xenv/test_path_remove.go

### Implementation for User Story 3

- [x] T044 [P] [US3] Implement environment variable management in pkg/xenv/env/manager.go
- [x] T045 [US3] Implement PATH management in pkg/xenv/env/path_manager.go
- [x] T046 [US3] Create xenv env command in cmd/kite/xenv_env_cmd.go
- [x] T047 [US3] Create xenv path command in cmd/kite/xenv_path_cmd.go
- [x] T048 [US3] Add env --set/--unset functionality to env command
- [x] T049 [US3] Add path --add/--rm functionality to path command
- [x] T050 [US3] Update ActivityState to manage active env vars and paths

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: User Story 4 - 导入导出配置 (Priority: P3)

**Goal**: 实现不同机器间同步开发环境配置，以便快速设置新的开发环境

**Independent Test**: 开发者可以使用`kite xenv config --export`导出配置，然后在另一台机器上使用`kite xenv config --import`导入

### Tests for User Story 4 (OPTIONAL - only if tests requested) ⚠️

- [ ] T051 [P] [US4] Contract test for `kite xenv config --export zip` in tests/integration/xenv/test_config_export.go
- [ ] T052 [P] [US4] Contract test for `kite xenv config --import <path>` in tests/integration/xenv/test_config_import.go

### Implementation for User Story 4

- [x] T053 [P] [US4] Implement configuration export functionality in pkg/xenv/config/exporter.go
- [x] T054 [US4] Implement configuration import functionality in pkg/xenv/config/importer.go
- [x] T055 [US4] Create xenv config command in cmd/kite/xenv_config_cmd.go
- [x] T056 [US4] Add config --export/--import functionality to config command
- [x] T057 [US4] Add config --set/--get functionality to config command
- [x] T058 [US4] Add validation for size limits (10MB) to export functionality

**Checkpoint**: All user stories should now be independently functional

---

## Phase 7: User Story 5 - Shell集成和实时生效 (Priority: P2)

**Goal**: 实现切换工具链后立即生效，无需重启shell或执行额外命令

**Independent Test**: 开发者配置shell hook `kite xenv shell --type bash` 后，切换工具版本应立即在当前shell中生效

### Tests for User Story 5 (OPTIONAL - only if tests requested) ⚠️

- [ ] T059 [P] [US5] Contract test for `kite xenv shell --type bash` in tests/integration/xenv/test_shell_bash.go
- [ ] T060 [P] [US5] Contract test for `kite xenv shell --type zsh` in tests/integration/xenv/test_shell_zsh.go
- [ ] T061 [P] [US5] Contract test for `kite xenv shell --type pwsh` in tests/integration/xenv/test_shell_pwsh.go

### Implementation for User Story 5

- [x] T062 [P] [US5] Implement shell hook generation for bash in pkg/xenv/shell/bash_hook.go
- [x] T063 [P] [US5] Implement shell hook generation for zsh in pkg/xenv/shell/zsh_hook.go
- [x] T064 [P] [US5] Implement shell hook generation for PowerShell in pkg/xenv/shell/pwsh_hook.go
- [x] T065 [US5] Create xenv shell command in cmd/kite/xenv_shell_cmd.go
- [x] T066 [US5] Implement directory-level config file processing (.xenv.toml, .envrc)
- [x] T067 [US5] Create shell integration utilities in pkg/xenv/shell/utils.go

**Checkpoint**: All user stories should now be independently functional

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T068 [P] Documentation updates for all commands in docs/xenv/
- [x] T069 Implementation of logging based on research decisions in pkg/xenv/logger.go
- [x] T070 [P] Error handling implementation based on research decisions across all modules
- [x] T071 Add xenv list command to integrate tool, env, and path listing functionality
- [x] T072 [P] Implement xenv init command for initialization process
- [x] T073 Add file size validation for configuration export (10MB limit)
- [x] T074 [P] Add support for .xenv.toml file processing in project directories
- [x] T075 Implement uninstall behavior allowing user to decide to keep config (from clarifications)
- [x] T076 Run quickstart.md validation to ensure all workflows work

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1/US2 but should be independently testable
- **User Story 4 (P3)**: Can start after Foundational (Phase 2) - May integrate with previous stories but should be independently testable
- **User Story 5 (P2)**: Can start after Foundational (Phase 2) - May integrate with previous stories but should be independently testable

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all models for User Story 1 together:
Task: "Implement ToolChain service in pkg/xenv/tools/service.go"
Task: "Implement tool installation logic in pkg/xenv/tools/installer.go"
Task: "Implement tool download functionality in pkg/xenv/tools/downloader.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Story 4 → Test independently → Deploy/Demo
6. Add User Story 5 → Test independently → Deploy/Demo
7. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
   - Developer D: User Story 4
   - Developer E: User Story 5
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

