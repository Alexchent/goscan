package cmd

import (
	"fmt"
	"github.com/Alexchent/goscan/logic"
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
		count += logic.SearchFromRedisSet(CacheKey, path)
		res := fmt.Sprintf("本次扫描发现 %d 个文件", count)
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
