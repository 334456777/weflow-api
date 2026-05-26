package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
)

var (
	msgTalker  string
	msgLimit   int
	msgOffset  int
	msgStart   string
	msgEnd     string
	msgKeyword string
	msgChatLab bool
	msgMedia   bool
	msgImage   string
	msgVoice   string
	msgVideo   string
	msgEmoji   string
)

func parseOptBool(s string) *bool {
	switch s {
	case "1", "true", "yes", "on":
		b := true
		return &b
	case "0", "false", "no", "off":
		b := false
		return &b
	}
	return nil
}

var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "拉取指定会话的消息（§3）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if msgTalker == "" {
			return fmt.Errorf("--talker 必填")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		p := client.MessagesParams{
			Talker:  msgTalker,
			Limit:   msgLimit,
			Offset:  msgOffset,
			Start:   msgStart,
			End:     msgEnd,
			Keyword: msgKeyword,
			Media:   msgMedia,
			Image:   parseOptBool(msgImage),
			Voice:   parseOptBool(msgVoice),
			Video:   parseOptBool(msgVideo),
			Emoji:   parseOptBool(msgEmoji),
		}
		r := newRenderer()
		if msgChatLab {
			resp, err := c.MessagesChatLab(cmd.Context(), p)
			if err != nil {
				return err
			}
			return r.MessagesChatLab(resp)
		}
		resp, err := c.Messages(cmd.Context(), p)
		if err != nil {
			return err
		}
		return r.Messages(resp)
	},
}

func init() {
	f := messagesCmd.Flags()
	f.StringVar(&msgTalker, "talker", "", "会话 ID（私聊 wxid 或群 xxx@chatroom），必填")
	f.IntVar(&msgLimit, "limit", 0, "返回条数（默认 100，1~10000）")
	f.IntVar(&msgOffset, "offset", 0, "分页偏移")
	f.StringVar(&msgStart, "start", "", "开始时间（YYYYMMDD 或时间戳）")
	f.StringVar(&msgEnd, "end", "", "结束时间（YYYYMMDD 自动扩到当天 23:59:59，或直接传时间戳）")
	f.StringVar(&msgKeyword, "keyword", "", "关键词过滤")
	f.BoolVar(&msgChatLab, "chatlab", false, "返回 ChatLab 格式")
	f.BoolVar(&msgMedia, "media", false, "导出并返回媒体地址")
	f.StringVar(&msgImage, "image", "", "media 启用时是否导出图片（1/0）")
	f.StringVar(&msgVoice, "voice", "", "media 启用时是否导出语音（1/0）")
	f.StringVar(&msgVideo, "video", "", "media 启用时是否导出视频（1/0）")
	f.StringVar(&msgEmoji, "emoji", "", "media 启用时是否导出表情（1/0）")
	_ = messagesCmd.MarkFlagRequired("talker")
	rootCmd.AddCommand(messagesCmd)
}
