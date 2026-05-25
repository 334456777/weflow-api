package cmd

import (
	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
)

var (
	contactsKeyword string
	contactsLimit   int
)

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "列出联系人（§5）",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := c.Contacts(cmd.Context(), client.ContactsParams{Keyword: contactsKeyword, Limit: contactsLimit})
		if err != nil {
			return err
		}
		return newRenderer().Contacts(resp)
	},
}

func init() {
	contactsCmd.Flags().StringVar(&contactsKeyword, "keyword", "", "按 username/nickname/remark/displayName 过滤")
	contactsCmd.Flags().IntVar(&contactsLimit, "limit", 0, "返回条数（默认 100）")
	rootCmd.AddCommand(contactsCmd)
}
