# dev


## Go程序包体积过大

### 使用构建标志优化

* `-s` 去除符号表
* `-w` 去除调试信息
* `-trimpath` 去除源码路径信息

```bash
# 基本优化构建
go build -ldflags="-s -w" -o app main.go

# 进一步优化
go build -ldflags="-s -w" -trimpath -o app main.go
```

### 使用工具分析包内容

可以使用以下工具来分析Go二进制文件的内容：

* go tool nm - 查看符号表
* go tool objdump - 反汇编查看代码
* go build -x - 查看构建过程
* 第三方工具 如 go-binsize 等

### 检查依赖包大小

```bash
# 查看依赖包大小
go list -f "{{.ImportPath}} {{.Size}}" ./...
```

* 使用 `go mod graph` 查看依赖关系
* 使用 `go list -m all` 查看所有依赖模块

### 常见解决方案

使用upx压缩：

```bash
go build -ldflags="-s -w" -o app
upx app
```
