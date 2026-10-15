// Command scan 是可 go install 安装的 scan 入口。
//
//	go install github.com/Alexchent/goscan/cmd/scan@latest
package main

import (
	"github.com/Alexchent/goscan/cmd"
)

func main() {
	cmd.Execute()
}
