package cmd

import (
	"fmt"
	"github.com/gookit/color"
	"regexp"
	"strconv"
	"strings"

	"github.com/Alexchent/goscan/cache"
	"github.com/Alexchent/goscan/help"
	"github.com/spf13/cobra"
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
	res := cache.SMembers(key)
	count = 0
	//过滤掉特殊字符-和_
	reg, _ := regexp.Compile("-|_")
	path = reg.ReplaceAllString(path, "")
	for _, val := range res {
		a := reg.ReplaceAllString(val, "")
		if strings.Contains(strings.ToLower(a), strings.ToLower(path)) {
			res := strings.Split(val, ",")
			if len(res) == 2 {
				fileSize, err := strconv.ParseInt(res[1], 10, 64)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println(res[0], color.HiGreen.Sprint(help.FormatFileSize(fileSize)))
				}
			} else {
				fmt.Println(val)
			}
			count++
		}
	}
	return
}
