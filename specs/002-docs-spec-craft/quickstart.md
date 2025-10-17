# Quickstart Guide: Kite XEnv

## 安装和初始化

首先，确保你已经安装了kite CLI工具。然后，按照以下步骤开始使用xenv：

```bash
# 检查xenv命令是否可用
kite xenv --help

# 初始化xenv配置
kite xenv init
```

## 基本使用

### 1. 管理开发工具

安装特定版本的开发工具：
```bash
# 安装Go 1.21版本
kite xenv tools install go@1.21

# 安装最新版本的Node.js
kite xenv tools install node@latest

# 查看所有已安装的工具
kite xenv list
```

### 2. 切换工具版本

在不同项目中切换不同版本的工具：
```bash
# 为当前会话激活Go 1.21版本
kite xenv use go@1.21

# 全局激活Node.js 18版本（在所有新会话中生效）
kite xenv use -g node@18

# 取消激活特定版本
kite xenv unuse go@1.21
```

### 3. 管理环境变量和PATH

设置环境变量：
```bash
# 为当前会话设置环境变量
kite xenv env --set NODE_ENV development

# 全局设置环境变量
kite xenv env --set -g GOPATH /path/to/gopath

# 查看所有环境变量
kite xenv env

# 删除环境变量
kite xenv env --unset NODE_ENV
```

添加和管理PATH路径：
```bash
# 添加路径到当前会话的PATH
kite xenv path --add ~/.my-tools/bin

# 全局添加路径
kite xenv path --add -g ~/.global-tools/bin

# 查看所有PATH条目
kite xenv path

# 从PATH中移除路径
kite xenv path --rm ~/.my-tools/bin
```

### 4. Shell集成

为了让环境切换立即生效，你需要在shell配置文件中添加hook：

对于Bash:
```bash
echo 'eval "$(kite xenv shell --type bash)"' >> ~/.bashrc
```

对于Zsh:
```bash
echo 'eval "$(kite xenv shell --type zsh)"' >> ~/.zshrc
```

对于PowerShell:
```powershell
echo 'kite xenv shell --type pwsh | Out-String | Invoke-Expression' >> $PROFILE
```

### 5. 配置管理

导出和导入配置：
```bash
# 导出当前配置为ZIP文件
kite xenv config --export zip

# 从ZIP文件导入配置
kite xenv config --import /path/to/config.zip

# 查看当前配置
kite xenv config

# 设置配置项
kite xenv config --set bin_dir ~/.local/bin
```

## 高级功能

### 项目级配置

在项目目录下创建`.xenv.toml`文件来配置项目的特定环境：
```toml
[tools]
node = "18.17.0"
go = "1.21.0"

[env]
NODE_ENV = "development"
GO_ENV = "dev"
```

### 管理多个工具版本

同时激活多个工具版本：
```bash
# 为当前会话激活多个工具版本
kite xenv use go@1.21 node@18.17

# 查看当前激活的工具
kite xenv list --activity
```

## 常见问题

### 工具安装失败

如果工具安装失败，请检查：
1. 网络连接是否正常
2. 目标版本是否存在
3. 目标目录是否有写入权限

### 环境变量未生效

如果环境变量在新shell中未生效，请检查：
1. 是否正确添加了shell hook
2. 是否使用了`-g`参数进行全局设置

### 需要更多帮助

使用`--help`参数查看各命令的帮助信息：
```bash
kite xenv --help
kite xenv tools --help
kite xenv use --help
```