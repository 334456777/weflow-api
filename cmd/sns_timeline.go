package cmd

import (
	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
)

var (
	snsTimelineLimit     int
	snsTimelineOffset    int
	snsTimelineUsernames string
	snsTimelineKeyword   string
	snsTimelineStart     string
	snsTimelineEnd       string
	snsTimelineMedia     string
	snsTimelineReplace   string
	snsTimelineInline    string
)

func parseOptInt(s string) *int {
	switch s {
	case "1", "true", "yes", "on":
		v := 1
		return &v
	case "0", "false", "no", "off":
		v := 0
		return &v
	}
	return nil
}

var snsTimelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "获取朋友圈时间线（§7.1）",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		raw, err := c.SnsTimeline(cmd.Context(), client.SnsTimelineParams{
			Limit:     snsTimelineLimit,
			Offset:    snsTimelineOffset,
			Usernames: snsTimelineUsernames,
			Keyword:   snsTimelineKeyword,
			Start:     snsTimelineStart,
			End:       snsTimelineEnd,
			Media:     parseOptInt(snsTimelineMedia),
			Replace:   parseOptInt(snsTimelineReplace),
			Inline:    parseOptInt(snsTimelineInline),
		})
		if err != nil {
			return err
		}
		return newRenderer().Raw(raw)
	},
}

func init() {
	f := snsTimelineCmd.Flags()
	f.IntVar(&snsTimelineLimit, "limit", 0, "返回数量（默认 20，1~200）")
	f.IntVar(&snsTimelineOffset, "offset", 0, "偏移量")
	f.StringVar(&snsTimelineUsernames, "usernames", "", "发布者过滤，逗号分隔")
	f.StringVar(&snsTimelineKeyword, "keyword", "", "正文关键词过滤")
	f.StringVar(&snsTimelineStart, "start", "", "起始时间（YYYYMMDD 或时间戳）")
	f.StringVar(&snsTimelineEnd, "end", "", "结束时间（YYYYMMDD 或时间戳）")
	f.StringVar(&snsTimelineMedia, "media", "", "返回可访问媒体地址（1/0，缺省 1）")
	f.StringVar(&snsTimelineReplace, "replace", "", "用解析地址覆盖 media.url/thumb（1/0，缺省 1）")
	f.StringVar(&snsTimelineInline, "inline", "", "内联返回 data: URL（1/0，缺省 0）")
	snsCmd.AddCommand(snsTimelineCmd)
}
