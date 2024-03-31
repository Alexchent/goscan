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
	"github.com/Alexchent/goscan/cache/mredis"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "查找文件",
	Long:  `从记录中查找是否存在指定文件，返回文件路径`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		fmt.Printf("请输入要查询的文件:\n")
		_, err := fmt.Scan(&path)
		if err != nil {
			return
		}
		count := 0
		count += SearchFromRedisSet(CacheKey, path)
		res := fmt.Sprintf("本次扫描发现 %d 个文件", count)
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}

func SearchFromRedisSet(key, path string) (count int) {
	res := mredis.SMembers(key)
	count = 0
	//过滤掉特殊字符-和_
	reg, _ := regexp.Compile("-|_")
	path = reg.ReplaceAllString(path, "")
	for _, val := range res {
		a := reg.ReplaceAllString(val, "")
		if strings.Contains(strings.ToLower(a), strings.ToLower(path)) {
			fmt.Println(val)
			count++
		}
	}
	return
}
