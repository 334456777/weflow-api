package cmd

import "github.com/spf13/cobra"

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "调用 /api/v1/health 健康检查（§1）",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := c.Health(cmd.Context())
		if err != nil {
			return err
		}
		return newRenderer().Health(resp)
	},
}

func init() { rootCmd.AddCommand(healthCmd) }
