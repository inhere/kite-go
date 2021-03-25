# link https://github.com/humbug/box/blob/master/Makefile
#SHELL = /bin/sh
.DEFAULT_GOAL := help
# 每行命令之前必须有一个tab键。如果想用其他键，可以用内置变量.RECIPEPREFIX 声明
# mac 下这条声明 没起作用 !!
.RECIPEPREFIX = >
.PHONY: all usage help clean

# 需要注意的是，每行命令在一个单独的shell中执行。这些Shell之间没有继承关系。
# - 解决办法是将两行命令写在一行，中间用分号分隔。
# - 或者在换行符前加反斜杠转义 \

# 接收命令行传入参数 make COMMAND tag=v2.0.4
# TAG=$(tag)

# 定义变量 使用 $(VAR1)
# VAR1=val

##There some make command for the project
##

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//' | sed -e 's/: / /'

##Available Commands:

  kit2gobin:     ## build kit to go bin dir
kit2gobin:
	go build -ldflags="-X 'kite.Info.Version=v2.0.2' -X 'kite.PubDate=222233'" -o $(GOPATH)/bin/kit ./bin/kit
	chmod a+x $(GOPATH)/bin/kit

  kite2gobin:     ## build kite to go bin dir
kite2gobin:
	go build  -o $(GOPATH)/bin/kit ./bin/kit
	chmod a+x $(GOPATH)/bin/kit

  csfix:      ## Fix code style for all files by go fmt
csfix:
	go fmt ./...

  test1:     ## Display code style error files by gofmt
test1:
	date
