package main

import (
	"fmt"
	"os"
)

// case 1
//
// windows git bash 运行: repl -f=- -t=/
// 得到的却是： win-bash.exe repl -f=- -t=C:/Program Files/Git/ ck-cz]
func main() {
	args := os.Args
	fmt.Println(args)
}
