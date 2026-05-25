package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/334456777/weflow-api/internal/client"
	"github.com/334456777/weflow-api/internal/config"
)

type TableRenderer struct {
	w     io.Writer
	color bool
}

func NewTable(w io.Writer, color bool) *TableRenderer {
	return &TableRenderer{w: w, color: color}
}

func (r *TableRenderer) newTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(r.w)
	t.SetStyle(table.StyleRounded)
	if !r.color {
		t.Style().Color = table.ColorOptions{}
	}
	return t
}

// ===== Health =====
func (r *TableRenderer) Health(v *client.HealthResponse) error {
	_, err := fmt.Fprintln(r.w, "状态:", v.Status)
	return err
}

// ===== Sessions =====
func (r *TableRenderer) Sessions(v *client.SessionsResponse) error {
	t := r.newTable()
	t.AppendHeader(table.Row{"#", "username", "displayName", "type", "lastTime", "unread"})
	for i, s := range v.Sessions {
		t.AppendRow(table.Row{i + 1, Truncate(s.Username, 24), Truncate(s.DisplayName, 20), s.Type, UnixSec(s.LastTimestamp), s.UnreadCount})
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 6, Align: text.AlignRight},
	})
	t.Render()
	fmt.Fprintf(r.w, "共 %d 个会话\n", v.Count)
	return nil
}

func (r *TableRenderer) SessionsChatLab(v *client.ChatLabSessionsResponse) error {
	t := r.newTable()
	t.AppendHeader(table.Row{"#", "id", "name", "type", "platform", "messageCount", "lastMessageAt"})
	for i, s := range v.Sessions {
		t.AppendRow(table.Row{i + 1, Truncate(s.ID, 24), Truncate(s.Name, 20), s.Type, s.Platform, s.MessageCount, UnixSec(s.LastMessageAt)})
	}
	t.Render()
	fmt.Fprintf(r.w, "共 %d 个会话 (ChatLab 格式)\n", len(v.Sessions))
	return nil
}

// ===== Messages =====
func (r *TableRenderer) Messages(v *client.MessagesResponse) error {
	t := r.newTable()
	t.AppendHeader(table.Row{"#", "time", "↑↓", "sender", "type", "content"})
	for i, m := range v.Messages {
		content := m.ParsedContent
		if content == "" {
			content = m.Content
		}
		t.AppendRow(table.Row{
			i + 1,
			UnixSec(m.CreateTime),
			SendArrow(m.IsSend),
			Truncate(m.SenderUsername, 18),
			m.LocalType,
			Truncate(content, 50),
		})
	}
	t.Render()
	fmt.Fprintf(r.w, "talker=%s  共 %d 条  hasMore=%s", v.Talker, v.Count, Boolyn(v.HasMore))
	if v.Media.Enabled {
		fmt.Fprintf(r.w, "  media=%d (exportPath=%s)", v.Media.Count, v.Media.ExportPath)
	}
	fmt.Fprintln(r.w)
	return nil
}

func (r *TableRenderer) MessagesChatLab(v *client.ChatLabMessagesResponse) error {
	fmt.Fprintf(r.w, "ChatLab v%s  generator=%s  exportedAt=%s\n",
		v.ChatLab.Version, v.ChatLab.Generator, UnixSec(v.ChatLab.ExportedAt))
	fmt.Fprintf(r.w, "Meta: name=%s  platform=%s  type=%s\n", v.Meta.Name, v.Meta.Platform, v.Meta.Type)
	if len(v.Members) > 0 {
		fmt.Fprintf(r.w, "成员 %d 个（已省略，加 --json 查看）\n", len(v.Members))
	}
	t := r.newTable()
	t.AppendHeader(table.Row{"#", "time", "sender", "type", "content"})
	for i, m := range v.Messages {
		name := m.AccountName
		if m.GroupNickname != "" {
			name = m.GroupNickname + " / " + name
		}
		t.AppendRow(table.Row{i + 1, UnixSec(m.Timestamp), Truncate(name, 20), m.Type, Truncate(m.Content, 50)})
	}
	t.Render()
	fmt.Fprintf(r.w, "共 %d 条消息\n", len(v.Messages))
	if v.Sync != nil {
		fmt.Fprintf(r.w, "Sync: hasMore=%s  nextSince=%d  nextOffset=%d  watermark=%s\n",
			Boolyn(v.Sync.HasMore), v.Sync.NextSince, v.Sync.NextOffset, UnixSec(v.Sync.Watermark))
	}
	return nil
}

// ===== Contacts =====
func (r *TableRenderer) Contacts(v *client.ContactsResponse) error {
	t := r.newTable()
	t.AppendHeader(table.Row{"#", "username", "displayName", "nickname", "remark", "alias", "type"})
	for i, c := range v.Contacts {
		t.AppendRow(table.Row{i + 1, Truncate(c.Username, 24), Truncate(c.DisplayName, 16), Truncate(c.Nickname, 16), Truncate(c.Remark, 12), Truncate(c.Alias, 12), c.Type})
	}
	t.Render()
	fmt.Fprintf(r.w, "共 %d 个联系人\n", v.Count)
	return nil
}

// ===== Members =====
func (r *TableRenderer) Members(v *client.GroupMembersResponse) error {
	t := r.newTable()
	t.AppendHeader(table.Row{"#", "wxid", "displayName", "groupNickname", "alias", "owner", "friend", "msgs"})
	for i, m := range v.Members {
		t.AppendRow(table.Row{i + 1, Truncate(m.WxID, 22), Truncate(m.DisplayName, 16), Truncate(m.GroupNickname, 12), Truncate(m.Alias, 12), Boolyn(m.IsOwner), Boolyn(m.IsFriend), m.MessageCount})
	}
	t.Render()
	fmt.Fprintf(r.w, "chatroom=%s  共 %d 人  fromCache=%s  updatedAt=%s\n",
		v.ChatroomID, v.Count, Boolyn(v.FromCache), UnixMillisOrSec(v.UpdatedAt))
	return nil
}

// ===== Event (SSE) =====
func (r *TableRenderer) Event(e client.Event) error {
	prefix := e.Name
	if r.color {
		switch e.Name {
		case "message.new":
			prefix = text.FgGreen.Sprint(e.Name)
		case "message.revoke":
			prefix = text.FgRed.Sprint(e.Name)
		default:
			prefix = text.FgYellow.Sprint(e.Name)
		}
	}
	if e.Data == "" {
		_, err := fmt.Fprintf(r.w, "[%s]\n", prefix)
		return err
	}
	// 尝试 pretty 解析 data（如果是 JSON）
	var v any
	if err := json.Unmarshal([]byte(e.Data), &v); err == nil {
		if m, ok := v.(map[string]any); ok {
			source, _ := m["sourceName"].(string)
			group, _ := m["groupName"].(string)
			content, _ := m["content"].(string)
			ts, _ := m["timestamp"].(float64)
			tsStr := "-"
			if ts > 0 {
				tsStr = UnixSec(int64(ts))
			}
			if group != "" {
				_, err := fmt.Fprintf(r.w, "[%s] %s  %s @ %s: %s\n", prefix, tsStr, source, group, content)
				return err
			}
			_, err := fmt.Fprintf(r.w, "[%s] %s  %s: %s\n", prefix, tsStr, source, content)
			return err
		}
	}
	_, err := fmt.Fprintf(r.w, "[%s] %s\n", prefix, e.Data)
	return err
}

// ===== Config =====
func (r *TableRenderer) Config(c *config.Config) error {
	t := r.newTable()
	t.AppendHeader(table.Row{"key", "value"})
	for _, kv := range ConfigSummary(c) {
		t.AppendRow(table.Row{kv[0], kv[1]})
	}
	t.Render()
	return nil
}

// ===== Raw (朋友圈系列) =====
func (r *TableRenderer) Raw(b json.RawMessage) error {
	s, err := IndentJSON(b)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(r.w, s)
	return err
}
