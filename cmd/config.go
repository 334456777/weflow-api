package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/config"
)

var (
	configInitForce bool
	configInitPath  string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看与初始化 CLI 配置",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "在默认位置写入空白配置模板",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := configInitPath
		if path == "" {
			path = config.DefaultPath()
		}
		if _, err := os.Stat(path); err == nil && !configInitForce {
			return fmt.Errorf("%s 已存在；加 --force 覆盖", path)
		}
		if err := config.WriteTemplate(path); err != nil {
			return err
		}
		fmt.Printf("已写入配置模板: %s\n", path)
		fmt.Println("请编辑该文件，把 token 改为你在 WeFlow 应用设置里复制的 Access Token。")
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "打印当前生效的配置（token 脱敏）",
	RunE: func(cmd *cobra.Command, args []string) error {
		return newRenderer().Config(cfg)
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "打印当前生效的配置文件路径",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.SourcePath == "" {
			fmt.Println("(没有加载任何配置文件；使用默认值/环境变量/flag)")
			return nil
		}
		fmt.Println(cfg.SourcePath)
		return nil
	},
}

func init() {
	configInitCmd.Flags().BoolVar(&configInitForce, "force", false, "已存在时覆盖写入")
	configInitCmd.Flags().StringVar(&configInitPath, "path", "", "指定写入路径（默认按 XDG）")
	configCmd.AddCommand(configInitCmd, configShowCmd, configPathCmd)
	rootCmd.AddCommand(configCmd)
}
