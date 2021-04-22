# something

gfwlist address:

- https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt

write custom rules:

```text
! Put user rules line by line in this file.
! See https://adblockplus.org/en/filter-cheatsheet
||en.wikipedia.org
||github.com
```

## pac resource

- https://github.com/jackwakefield/gopac/blob/master/runtime.go

### gfwlist2pac

https://github.com/petronny/gfwlist2pac

- https://raw.githubusercontent.com/petronny/gfwlist2pac/master/gfwlist.pac

### pac file

```js
var proxy = 'SOCKS5 127.0.0.1:1080';
// multi
var proxy = "SOCKS5 127.0.0.1:9527; PROXY 127.0.0.1:9580; DIRECT"
```

## refer

- 代理自动配置 https://zh.wikipedia.org/wiki/%E4%BB%A3%E7%90%86%E8%87%AA%E5%8A%A8%E9%85%8D%E7%BD%AE
- https://zh.wikipedia.org/zh/%E4%BB%A3%E7%90%86%E8%87%AA%E5%8A%A8%E9%85%8D%E7%BD%AE
