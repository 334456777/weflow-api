package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/334456777/weflow-api/internal/client"
	"github.com/334456777/weflow-api/internal/config"
)

// NameStyle 控制 --text 模式下成员名的显示风格。
type NameStyle int

const (
	NameStyleFull     NameStyle = iota // "群昵称(本名)"；无群昵称时 fallback 到本名
	NameStyleNickname                  // 只群昵称；无群昵称时 fallback 到本名
	NameStyleAccount                   // 只本名
)

func ParseNameStyle(s string) NameStyle {
	switch s {
	case "nickname":
		return NameStyleNickname
	case "account":
		return NameStyleAccount
	default:
		return NameStyleFull
	}
}

type TextRenderer struct {
	w         io.Writer
	nameStyle NameStyle
}

func NewText(w io.Writer, nameStyle NameStyle) *TextRenderer {
	return &TextRenderer{w: w, nameStyle: nameStyle}
}

func (r *TextRenderer) formatName(account, nick string) string {
	switch r.nameStyle {
	case NameStyleNickname:
		if nick != "" {
			return nick
		}
		return account
	case NameStyleAccount:
		return account
	default:
		if nick != "" {
			return nick + "(" + account + ")"
		}
		return account
	}
}

func (r *TextRenderer) Health(v *client.HealthResponse) error {
	_, err := fmt.Fprintln(r.w, v.Status)
	return err
}

func (r *TextRenderer) Sessions(v *client.SessionsResponse) error {
	for _, s := range v.Sessions {
		fmt.Fprintf(r.w, "%s\t%s\ttype=%d\tlast=%s\tunread=%d\n",
			s.Username, s.DisplayName, s.Type, UnixSec(s.LastTimestamp), s.UnreadCount)
	}
	return nil
}

func (r *TextRenderer) SessionsChatLab(v *client.ChatLabSessionsResponse) error {
	for _, s := range v.Sessions {
		fmt.Fprintf(r.w, "%s\t%s\t%s\tmsgs=%d\tlast=%s\n",
			s.ID, s.Name, s.Type, s.MessageCount, UnixSec(s.LastMessageAt))
	}
	return nil
}

func (r *TextRenderer) Messages(v *client.MessagesResponse) error {
	for _, m := range v.Messages {
		content := m.ParsedContent
		if content == "" {
			content = m.Content
		}
		who := m.SenderUsername
		if m.IsSend == 1 {
			who = "[我] " + who
		}
		fmt.Fprintf(r.w, "[%s] %s: %s\n", UnixSec(m.CreateTime), who, FormatContent(content))
	}
	return nil
}

func (r *TextRenderer) MessagesChatLab(v *client.ChatLabMessagesResponse) error {
	for _, m := range v.Messages {
		fmt.Fprintf(r.w, "[%s] %s: %s\n", UnixSec(m.Timestamp), r.formatName(m.AccountName, m.GroupNickname), FormatContent(m.Content))
	}
	return nil
}

func (r *TextRenderer) Contacts(v *client.ContactsResponse) error {
	for _, c := range v.Contacts {
		extra := ""
		if c.Remark != "" {
			extra = " (备注: " + c.Remark + ")"
		}
		fmt.Fprintf(r.w, "%s\t%s%s\n", c.Username, c.DisplayName, extra)
	}
	return nil
}

func (r *TextRenderer) Members(v *client.GroupMembersResponse) error {
	for _, m := range v.Members {
		suffix := ""
		if m.IsOwner {
			suffix = " [群主]"
		}
		if m.MessageCount > 0 {
			suffix += fmt.Sprintf(" msgs=%d", m.MessageCount)
		}
		fmt.Fprintf(r.w, "%s\t%s%s\n", m.WxID, r.formatName(m.DisplayName, m.GroupNickname), suffix)
	}
	return nil
}

func (r *TextRenderer) Event(e client.Event) error {
	if e.Data == "" {
		_, err := fmt.Fprintf(r.w, "[%s]\n", e.Name)
		return err
	}
	var v any
	if err := json.Unmarshal([]byte(e.Data), &v); err == nil {
		if m, ok := v.(map[string]any); ok {
			source, _ := m["sourceName"].(string)
			group, _ := m["groupName"].(string)
			content, _ := m["content"].(string)
			content = FormatContent(content)
			tsStr := "-"
			if ts, ok := m["timestamp"].(float64); ok && ts > 0 {
				tsStr = UnixSec(int64(ts))
			}
			if group != "" {
				_, err := fmt.Fprintf(r.w, "[%s] %s %s @ %s: %s\n", e.Name, tsStr, source, group, content)
				return err
			}
			_, err := fmt.Fprintf(r.w, "[%s] %s %s: %s\n", e.Name, tsStr, source, content)
			return err
		}
	}
	_, err := fmt.Fprintf(r.w, "[%s] %s\n", e.Name, e.Data)
	return err
}

func (r *TextRenderer) Config(c *config.Config) error {
	for _, kv := range ConfigSummary(c) {
		fmt.Fprintf(r.w, "%s: %s\n", kv[0], kv[1])
	}
	return nil
}

func (r *TextRenderer) Raw(b json.RawMessage) error {
	// 朋友圈类响应字段不可知，没有简单 text 模板，回退到 pretty JSON
	s, err := IndentJSON(b)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(r.w, s)
	return err
}
