package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var snsBlockDeleteCmd = &cobra.Command{
	Use:   "block-delete <status|install|uninstall>",
	Short: "朋友圈防删开关（§7.6）",
	Args:  cobra.ExactArgs(1),
	ValidArgs: []string{"status", "install", "uninstall"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		var raw json.RawMessage
		switch args[0] {
		case "status":
			raw, err = c.SnsBlockDeleteStatus(ctx)
		case "install":
			raw, err = c.SnsBlockDeleteInstall(ctx)
		case "uninstall":
			raw, err = c.SnsBlockDeleteUninstall(ctx)
		default:
			return fmt.Errorf("未知子动作 %q，期望 status/install/uninstall", args[0])
		}
		if err != nil {
			return err
		}
		return newRenderer().Raw(raw)
	},
}

func init() { snsCmd.AddCommand(snsBlockDeleteCmd) }
