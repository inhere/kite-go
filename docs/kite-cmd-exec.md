# kite 命令执行原理

- 先判断是否是内置命令 是内置命令，则执行内置命令
- 外部扩展命令搜索。判断是否是 `kite-<ext>` 扩展命令
  - 在系统 PATH 环境变量指定的目录中查找名为 `kite-<ext>` 的可执行文件。
- 别名（Alias）解析 - 检查是否为用户定义的别名 `aliases` 配置
- kscript 脚本解析 - 检查是否为 kscript 脚本, task 定义等


### kite ext

参考 git 扩展的实现机制，扩展命令（自定义子命令）的加载实现原理主要基于 命令查找机制 和 可执行文件命名约定。

以下是详细解析：

- 需要先添加 extension `kite app ext add <ext-name> <ext-description>`
  - 会在系统 PATH 环境变量指定的目录中查找名为 `kite-<ext-name>` 的可执行文件
  - 运行该可执行文件，校验 ext 存在且可执行
- 后续就可以使用 `kite <ext-name>` 运行该扩展命令
