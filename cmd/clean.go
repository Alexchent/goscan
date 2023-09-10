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
	"github.com/Alexchent/goscan/cache/redis"
	"strings"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理掉制定的文件类型",
	Long:  `清理掉制定的文件类型`,
	Run: func(cmd *cobra.Command, args []string) {
		key := "have_save_file"
		val := redis.SMembers(key)
		for _, v := range val {
			newv := strings.TrimRight(v, "\n")
			redis.SRem(key, v)
			redis.SAdd(key, newv)

			// 按 后缀清理
			if strings.HasSuffix(v, "js") ||
				strings.HasSuffix(v, "torrent") ||
				strings.HasSuffix(v, "jpeg") {
				fmt.Println("过滤掉:", v)
				redis.SRem(key, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}