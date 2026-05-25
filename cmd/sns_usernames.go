package cmd

import "github.com/spf13/cobra"

var snsUsernamesCmd = &cobra.Command{
	Use:   "usernames",
	Short: "获取朋友圈发布者列表（§7.2）",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		raw, err := c.SnsUsernames(cmd.Context())
		if err != nil {
			return err
		}
		return newRenderer().Raw(raw)
	},
}

func init() { snsCmd.AddCommand(snsUsernamesCmd) }
