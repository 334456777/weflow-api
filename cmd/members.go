package cmd

import (
	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
)

var (
	membersWithCounts   bool
	membersForceRefresh bool
)

var membersCmd = &cobra.Command{
	Use:   "members <chatroomId>",
	Short: "列出群成员（§6）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := c.Members(cmd.Context(), client.MembersParams{
			ChatroomID:   args[0],
			WithCounts:   membersWithCounts,
			ForceRefresh: membersForceRefresh,
		})
		if err != nil {
			return err
		}
		return newRenderer().Members(resp)
	},
}

func init() {
	membersCmd.Flags().BoolVar(&membersWithCounts, "with-counts", false, "附带每位成员发言数")
	membersCmd.Flags().BoolVar(&membersForceRefresh, "force-refresh", false, "跳过内存缓存强制刷新")
	rootCmd.AddCommand(membersCmd)
}
