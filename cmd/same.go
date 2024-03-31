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
	fmt "fmt"
	scan "github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
	"time"
)

// startCmd represents the start command
var sameCmd = &cobra.Command{
	Use:   "same",
	Short: "查找重复文件",
	Long:  `查找重复文件`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("扫描完成，耗时：", time.Since(start))

		if path == "" {
			fmt.Printf("请输入要扫描的目录:\n")
			_, _ = fmt.Scan(&path)
		}

		scan.Search(path)
		// 输出重复文件到txt中
		//scan.GetSame()
		fmt.Println("记录扫描文件 start")
		same := scan.LogSameFile()
		fmt.Println("记录扫描文件 end")
		if same == nil {
			return
		}
		fmt.Println("执行删除操作 start")
		scan.RemoveSameFile(same)
		fmt.Println("执行删除操作 end")
	},
}

func init() {
	sameCmd.Flags().StringVarP(&path, "dir", "d", "", "需要扫描的目录")
	rootCmd.AddCommand(sameCmd)
}
