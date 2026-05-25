package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	snsProxyURL string
	snsProxyKey string
	snsProxyOut string
)

var snsProxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "朋友圈媒体代理下载（§7.4）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if snsProxyURL == "" {
			return fmt.Errorf("--url 必填")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		var dst io.Writer = os.Stdout
		if snsProxyOut != "" {
			f, err := os.Create(snsProxyOut)
			if err != nil {
				return err
			}
			defer f.Close()
			dst = f
		}
		ct, err := c.SnsProxy(cmd.Context(), snsProxyURL, snsProxyKey, dst)
		if err != nil {
			return err
		}
		if snsProxyOut != "" {
			fmt.Fprintf(os.Stderr, "已下载到 %s (Content-Type: %s)\n", snsProxyOut, ct)
		}
		return nil
	},
}

func init() {
	snsProxyCmd.Flags().StringVar(&snsProxyURL, "url", "", "媒体原始 URL，必填")
	snsProxyCmd.Flags().StringVar(&snsProxyKey, "key", "", "解密 key（部分资源需要）")
	snsProxyCmd.Flags().StringVarP(&snsProxyOut, "output", "o", "", "输出文件（缺省 stdout）")
	_ = snsProxyCmd.MarkFlagRequired("url")
	snsCmd.AddCommand(snsProxyCmd)
}
