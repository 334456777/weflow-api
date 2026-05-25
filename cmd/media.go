package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var mediaOut string

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "访问导出媒体（§8）",
}

var mediaGetCmd = &cobra.Command{
	Use:   "get <relativePath>",
	Short: "下载一个媒体文件，缺省 -o 时输出到 stdout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		var dst io.Writer = os.Stdout
		if mediaOut != "" {
			f, err := os.Create(mediaOut)
			if err != nil {
				return err
			}
			defer f.Close()
			dst = f
		}
		ct, err := c.MediaGet(cmd.Context(), args[0], dst)
		if err != nil {
			return err
		}
		if mediaOut != "" {
			fmt.Fprintf(os.Stderr, "已下载到 %s (Content-Type: %s)\n", mediaOut, ct)
		}
		return nil
	},
}

func init() {
	mediaGetCmd.Flags().StringVarP(&mediaOut, "output", "o", "", "输出文件路径（缺省 stdout）")
	mediaCmd.AddCommand(mediaGetCmd)
	rootCmd.AddCommand(mediaCmd)
}
