package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
	"github.com/334456777/weflow-api/internal/config"
	"github.com/334456777/weflow-api/internal/output"
)

var (
	cfgFile     string
	flagHost    string
	flagToken   string
	flagTimeout time.Duration
	flagJSON      bool
	flagText      bool
	flagNoColor   bool
	flagNameStyle string

	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:           "weflow",
	Short:         "WeFlow API 命令行客户端",
	Long:          "调用本机 WeFlow HTTP / SSE API 的命令行工具。先 `weflow config init` 生成配置文件，填入 token 与 host 即可使用。",
	SilenceUsage:  true,
	SilenceErrors: false,
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径（默认依次搜索 ./weflow.yaml、$XDG_CONFIG_HOME/weflow/config.yaml、~/.weflow.yaml）")
	rootCmd.PersistentFlags().StringVar(&flagHost, "host", "", "覆盖 WeFlow 基础地址，如 http://127.0.0.1:5031")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "覆盖 Access Token")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 0, "HTTP 超时（默认 30s）")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "以原样 JSON 输出")
	rootCmd.PersistentFlags().BoolVar(&flagText, "text", false, "纯文本输出（每行一条记录，适合管道/重定向到 .txt）")
	rootCmd.PersistentFlags().StringVar(&flagNameStyle, "name-style", "full", "--text 模式下成员名风格：full（群昵称(本名)）| nickname（只群昵称）| account（只本名）")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "禁用 ANSI 颜色")
}

func initConfig() {
	c, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	flags := rootCmd.PersistentFlags()
	if flags.Changed("host") {
		c.Host = flagHost
	}
	if flags.Changed("token") {
		c.Token = flagToken
	}
	if flags.Changed("timeout") {
		c.Timeout = flagTimeout
	}
	if flagNoColor {
		c.Output.Color = false
	}

	cfg = c
}

func newClient() (*client.Client, error) {
	return client.New(cfg)
}

func newRenderer() output.Renderer {
	if flagJSON {
		return output.NewJSON(os.Stdout)
	}
	if flagText {
		return output.NewText(os.Stdout, output.ParseNameStyle(flagNameStyle))
	}
	return output.NewTable(os.Stdout, cfg.Output.Color)
}
