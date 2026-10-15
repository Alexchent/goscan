package cmd

import (
	"fmt"
	"time"

	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/pkg"
	"github.com/Alexchent/goscan/service"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "开始扫描",
	Long:  `开始扫描`,
	Run: func(cmd *cobra.Command, args []string) {
		startWork()
	},
}

func init() {
	startCmd.Flags().StringVarP(&scanDir, "dir", "d", "", "需要扫描的目录")
	rootCmd.AddCommand(startCmd)
}

func startWork() {
	start := time.Now()
	defer func() {
		fmt.Println("扫描完成，耗时：", time.Since(start))
	}()

	// 获取扫描路径
	scanPath, err := pkg.GetScanPath(scanDir)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	// 验证路径
	if err := pkg.ValidatePath(scanPath); err != nil {
		fmt.Printf("路径验证失败: %v\n", err)
		return
	}

	color.HiGreen.Println("开始扫描：", scanPath)
	if err := service.Save(scanPath, config.FilterSuffix); err != nil {
		fmt.Printf("扫描出错: %v\n", err)
	}
}
