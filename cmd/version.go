package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("weflow %s\n  commit: %s\n  date:   %s\n", version.Version, version.Commit, version.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
