# proxy server

简单的http代理服务, 用于解决本地开发调试请求其他环境服务的问题.

### 功能说明

- 代理服务会将请求转发到指定的目标服务, 并将响应返回给客户端.
- 支持将请求转发到本地文件系统, 用于mock api响应数据.
- 支持修改请求数据, 如header, body等.
- 支付修改响应数据, 如header, body等.
- 支持使用正则表达式匹配请求.
- 支持使用变量引用上下文.
- 支持使用过滤器改变请求, 响应数据.

## 使用

```bash
$ kite dev proxy-server start
```
## 命令选项

```bash
$ kite dev proxy-server start --help
```

## 配置使用

每行一个规则,规则格式：`pattern distOperateURI`

- `pattern` 为匹配的请求host,路径等, 支持正则表达式.
- `distOperateURI` 为目标资源/服务的地址, 支持变量.

### 支持协议

- `file://` 本地文件系统

### 可用上下文

- `req.uriPath` 请求的路径, 不包含host部分
- `req.fullUri` 请求的完整路径, 包含host部分
- `req.host` 请求的host部分, 不包含端口
- `req.port` 请求的端口
- `req.protocol` 请求的协议 eg: http, https

### 可用的变量

## 配置示例

### 返回本地文件内容

```bash
# 参考了 https://wproxy.org/whistle 的规则设置
#
# 规则格式： pattern distOperateURI
# 增强功能：
# - distOperateURI 支持通过变量应用一些上下文

# 直接使用本地文件mock api响应数据
my-local-server:8080/api/some/detail file://C:/some/dir/api-mock/tmp/test-resp.json

# distOperateURI 里使用变量参考
# eg: 访问 my-local-server:8080/api/some/detail -> 会响应 C:/some/dir/api-mock/api/some/detail.json 的内容
my-local-server:8080 file://C:/some/dir/api-mock/{req.uriPath}.json

# 配置本地路径为静态资源目录
# pattern为域名或路径, 会自动根据请求url后面剩余的路径跟filepath自动补全
# eg: 访问 my-local-server:8080/api/some-detail.json -> 会响应 C:/some/dir/api-mock/api/some-detail.json 的内容
my-local-server:8080/api/ file://C:/some/dir/api-mock/api/

```

### 转发请求到其他服务

```bash
# 转发请求到其他服务
# eg: 访问 my-local-server:8080/api/some/detail -> 会转发到 http://other-server.com
my-local-server:8080 http://other-server.com

# 转发请求到其他服务, 并修改请求头
# eg: 访问 my-local-server:8080/api/some/detail -> 会转发到 http://other-server.com
# 并将请求的header中的`Host`字段替换为other-server.com
my-local-server:8080 http://other-server.com <<{
# add req header
add_header{
    Host: other-server.com
    New-Key: some-value
    Content-Type: application/json
}

set_body {
   {"key": "value"}
}

}>>

```
