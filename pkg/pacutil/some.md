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
