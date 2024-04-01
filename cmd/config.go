package cmd

import (
	_ "embed"
	"fmt"
	myFile "github.com/Alexchent/goscan/file"
	help "github.com/Alexchent/goscan/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
)

//go:embed scan.yaml
var configData string

var cover bool
var outputConf string

// cleanCmd represents the clean command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "生成默认配置文件",
	Long:  `生成默认配置文件`,
	Run: func(cmd *cobra.Command, args []string) {

		myFile.CreateDateDir(outputConf)
		filename := outputConf + "/scan.yaml"

		//fmt.Println(filename)
		//return
		if help.FileIsExist(filename) && cover == false {
			fmt.Println("配置文件已存在")
			return
		}

		fd, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		_, err := fd.WriteString(configData)
		if err != nil {
			panic("配置文件生成失败")
			return
		}
	},
}

// go run main clean -c apk
func init() {
	// 本地标志, 此处进队cleanCmd有效
	home, _ := homedir.Dir()
	//fmt.Println("home:", home)
	configCmd.Flags().BoolVar(&cover, "cover", false, "配置文件已存在的情况是否覆盖")
	configCmd.Flags().StringVarP(&outputConf, "output", "o", home, "配置文件输出路径")
	rootCmd.AddCommand(configCmd)
}
