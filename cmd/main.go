package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wenit/go-mod/internal/cmd/clean"
	"github.com/wenit/go-mod/internal/cmd/install"
	"github.com/wenit/go-mod/internal/cmd/pack"
	"github.com/wenit/go-mod/internal/version"
)

func main() {

	Execute()
}

// 参数
var (
	help  bool // 打印帮助信息
	ver   bool // 打印版本信息
	debug bool // 开启调试模式
)

// 根命令
var rootCmd = &cobra.Command{
	DisableFlagsInUseLine: true,
	Use:                   os.Args[0],
	Short:                 "使用说明",
	Long:                  `使用说明：`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if ver {
			fmt.Println(version.GetVersion())
			return
		}
	},
}

// Execute 执行程序
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().BoolVarP(&help, "help", "h", false, "帮助信息")
	rootCmd.Flags().BoolVarP(&ver, "version", "v", false, "版本信息")

	rootCmd.AddCommand(clean.GetSubCmd())
	rootCmd.AddCommand(pack.GetSubCmd())
	rootCmd.AddCommand(install.GetSubCmd())

	// 帮助文档
	rootCmd.SetHelpCommand(helpCmd)
}

var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "更多帮助文档",
	Long:  `Help provides help for any command in the application.`,

	Run: func(c *cobra.Command, args []string) {
		cmd, _, e := c.Root().Find(args)
		if cmd == nil || e != nil {
			c.Printf("Unknown help topic %#q\n", args)
			c.Root().Usage()
		} else {
			cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
			cmd.Help()
		}
	},
}
