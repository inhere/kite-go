package main

import (
	"context"

	"github.com/d5/tengo/v2"
	// "github.com/d5/tengo/v2/stdlib"
)

// see https://github.com/d5/tengo/blob/master/docs/stdlib-os.md
func main() {
	src := []byte(`
	os := import("os")
	fmt := import("fmt")

	fmt.println("hello world, from tengo")
	// 文件系统
	os.write("/tmp/tengo.txt", "hello from tengo")
	text := os.read("/tmp/tengo.txt")
	println("read back:", text)

	// 命令执行
	res := os.run("echo", "hello", "world")
	println("exec.ok =", res.ok)
	println("exec.out=", res.out)
	`)

	// 1. 创建脚本
	script := tengo.NewScript(src)

	// set values
	_ = script.Add("a", 1)

	// run the script
	compiled, err := script.RunContext(context.Background())
	if err != nil {
		panic(err)
	}
	_ = compiled.Get("a")
}
