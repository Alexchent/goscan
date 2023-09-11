package cmd

import (
	"fmt"
	mconf "github.com/Alexchent/goscan/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. This is goscan`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.0")

		//fmt.Println(viper.AllKeys())
		//fmt.Println("userLicense:", userLicense)
		//fmt.Println("author:", viper.Get("author"))
		//fmt.Println("useviper:", viper.Get("useviper"))
		//fmt.Println("license:", viper.Get("license"))
		fmt.Println(mconf.Conf.Dir)
		fmt.Println(mconf.Conf.FilterType)
	},
}
