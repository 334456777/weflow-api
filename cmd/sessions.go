package cmd

import (
	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
)

var (
	sessionsKeyword string
	sessionsLimit   int
	sessionsFormat  string
)

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "列出会话（§4 / §4.1）",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		p := client.SessionsParams{Keyword: sessionsKeyword, Limit: sessionsLimit}
		r := newRenderer()
		if sessionsFormat == "chatlab" {
			resp, err := c.SessionsChatLab(cmd.Context(), p)
			if err != nil {
				return err
			}
			return r.SessionsChatLab(resp)
		}
		resp, err := c.Sessions(cmd.Context(), p)
		if err != nil {
			return err
		}
		return r.Sessions(resp)
	},
}

func init() {
	sessionsCmd.Flags().StringVar(&sessionsKeyword, "keyword", "", "按 username/displayName 过滤")
	sessionsCmd.Flags().IntVar(&sessionsLimit, "limit", 0, "返回条数（默认 100）")
	sessionsCmd.Flags().StringVar(&sessionsFormat, "format", "", "传 chatlab 走 §4.1 ChatLab Pull 协议结构")
	rootCmd.AddCommand(sessionsCmd)
}
