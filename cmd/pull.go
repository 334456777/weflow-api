package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
)

var (
	pullSince  int64
	pullEnd    int64
	pullLimit  int
	pullOffset int
	pullFollow bool
)

var pullCmd = &cobra.Command{
	Use:   "pull <sessionId>",
	Short: "ChatLab Pull 拉会话消息（§4.2），可 --follow 自动分页",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		c, err := newClient()
		if err != nil {
			return err
		}
		r := newRenderer()
		p := client.PullParams{Since: pullSince, End: pullEnd, Limit: pullLimit, Offset: pullOffset}

		for page := 1; ; page++ {
			resp, err := c.Pull(cmd.Context(), sessionID, p)
			if err != nil {
				return err
			}
			if pullFollow {
				fmt.Fprintf(os.Stderr, "== 第 %d 页 ==\n", page)
			}
			if err := r.MessagesChatLab(resp); err != nil {
				return err
			}
			if !pullFollow || resp.Sync == nil || !resp.Sync.HasMore {
				return nil
			}
			p.Since = resp.Sync.NextSince
			p.Offset = resp.Sync.NextOffset
		}
	},
}

func init() {
	f := pullCmd.Flags()
	f.Int64Var(&pullSince, "since", 0, "起始秒级 Unix 时间戳，仅返回该时间之后的消息")
	f.Int64Var(&pullEnd, "end", 0, "结束秒级 Unix 时间戳")
	f.IntVar(&pullLimit, "limit", 0, "单页上限（默认且最大 5000）")
	f.IntVar(&pullOffset, "offset", 0, "分页偏移")
	f.BoolVar(&pullFollow, "follow", false, "自动翻页直到 hasMore=false")
	rootCmd.AddCommand(pullCmd)
}
