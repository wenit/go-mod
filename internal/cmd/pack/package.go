package pack

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/wenit/go-mod/pkg/pack"
)

var (
	ver      string // 输出版本号
	toDir    string // 输出文件
	fromDir  string // 输入文件
	excludes string // 排除目录，多个目录使用逗号分割，示例：.svn,.git,.vscode
)

// 子命令
var subCmd = &cobra.Command{
	Use:   "package",
	Short: "项目打包，用于发布",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(args)
	},
}

// GetSubCmd 获取子命令
func GetSubCmd() *cobra.Command {
	return subCmd
}

func run(args []string) error {
	path := fromDir
	version := ver

	outputDirectory := toDir

	if err := pack.Package(path, version, outputDirectory, excludes); err != nil {
		return fmt.Errorf("package module: %w", err)
	}
	return nil
}

func init() {
	flagSet := pflag.NewFlagSet("flag", pflag.ContinueOnError)
	flagSet.StringVarP(&fromDir, "from", "f", ".", "输入目录")
	flagSet.StringVarP(&toDir, "to", "t", "./target", "输出目录")
	flagSet.StringVarP(&ver, "version", "v", "v1.0.0", "输出版本号")
	flagSet.StringVarP(&excludes, "excludes", "e", ".svn,.git,.vscode,target,releases", "排除目录，多个目录使用逗号分割")
	subCmd.Flags().AddFlagSet(flagSet)
}
