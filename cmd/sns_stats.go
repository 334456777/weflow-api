package cmd

import (
	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
)

var snsStatsFast bool

var snsStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "朋友圈导出统计（§7.3）",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		raw, err := c.SnsExportStats(cmd.Context(), client.SnsExportStatsParams{Fast: snsStatsFast})
		if err != nil {
			return err
		}
		return newRenderer().Raw(raw)
	},
}

func init() {
	snsStatsCmd.Flags().BoolVar(&snsStatsFast, "fast", false, "走快速统计（优先缓存）")
	snsCmd.AddCommand(snsStatsCmd)
}
