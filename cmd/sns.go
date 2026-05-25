package cmd

import "github.com/spf13/cobra"

var snsCmd = &cobra.Command{
	Use:   "sns",
	Short: "朋友圈相关接口（§7）",
}

func init() { rootCmd.AddCommand(snsCmd) }
