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
	"encoding/json"
	"fmt"
	scan "github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var file string

// startCmd represents the start command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "移动重复文件",
	Long:  `移动重复文件`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("扫描完成，耗时：", time.Since(start))

		if file == "" {
			fmt.Println(file + "\t文件不存在")
			return
		}

		readFile, err := os.ReadFile(file)
		if err != nil {
			panic("打开文件失败: " + err.Error())
			return
		}

		var same scan.FileList
		err = json.Unmarshal(readFile, &same)
		if err != nil {
			log.Fatal("error during unmarshal():", err)
			return
		}
		//fmt.Println(same)
		fmt.Println("执行删除操作 start")
		scan.RemoveSameFile(same)
		fmt.Println("执行删除操作 end")
	},
}

func init() {
	removeCmd.Flags().StringVarP(&file, "file", "f", "", "same命令生成的文件")
	rootCmd.AddCommand(removeCmd)
}
