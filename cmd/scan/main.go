// Command scan 是 scan 的安装入口（也是唯一入口）。
//
//	go install github.com/Alexchent/goscan/cmd/scan@latest
//
// 二进制名为 scan。根目录不再保留重复的 main.go，构建/安装统一走这里。
package main

import (
	"github.com/Alexchent/goscan/cmd"
)

func main() {
	cmd.Execute()
}
