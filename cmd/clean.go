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
	"github.com/Alexchent/goscan/cache"
	"strings"

	"github.com/spf13/cobra"
)

var clearSuffix string
var clearContain string

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理掉符合条件的的文件",
	Long:  `清理掉符合条件的的文件`,
	Run: func(cmd *cobra.Command, args []string) {
		key := CacheKey
		val := cache.SMembers(key)
		for _, v := range val {
			// 按 后缀清理
			if clearSuffix != "" {
				if strings.HasSuffix(v, clearSuffix) {
					fmt.Println("过滤掉：", v)
					cache.SRem(key, v)
				}
			}

			// 按 后缀清理
			if clearContain != "" {
				if strings.Contains(v, clearContain) {
					fmt.Println("过滤掉：", v)
					cache.SRem(key, v)
				}
			}
		}
	},
}

// 后缀
// go run main clean -c apk
func init() {
	// 本地标志, 此处进队cleanCmd有效
	cleanCmd.Flags().StringVar(&clearSuffix, "suffix", "", "需要清理掉的文件类型")
	cleanCmd.Flags().StringVarP(&clearContain, "contain", "C", "", "包含该值")
	rootCmd.AddCommand(cleanCmd)
}
