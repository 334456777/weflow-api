package cmd

import "github.com/spf13/cobra"

var snsDeleteCmd = &cobra.Command{
	Use:   "delete <postId>",
	Short: "删除单条朋友圈（§7.7）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		raw, err := c.SnsDeletePost(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return newRenderer().Raw(raw)
	},
}

func init() { snsCmd.AddCommand(snsDeleteCmd) }
