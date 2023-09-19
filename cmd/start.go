/*
Copyright © 2023 Chen Tao <1023615292@qq.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	mconf "github.com/Alexchent/goscan/config"
	scan "github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "开始扫描",
	Long:  `开始扫描`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("扫描完成，耗时：", time.Since(start))

		fmt.Println("ignore:", mconf.Conf.FilterType)

		var path string
		fmt.Printf("请输入要扫描的目录:\n")

		_, err := fmt.Scan(&path)
		if err != nil {
			return
		}
		if path == "" || path == "/" {
			path = "/Users/chentao/Downloads"
		}
		scan.WriteToFile(path)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
