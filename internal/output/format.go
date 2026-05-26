package output

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/334456777/weflow-api/internal/config"
)

func UnixSec(ts int64) string {
	if ts <= 0 {
		return "-"
	}
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func UnixMillisOrSec(ts int64) string {
	if ts <= 0 {
		return "-"
	}
	if ts >= 1_000_000_000_000 {
		return time.UnixMilli(ts).Format("2006-01-02 15:04:05")
	}
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func Truncate(s string, n int) string {
	if n <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}

func MaskToken(s string) string {
	if s == "" {
		return "(空)"
	}
	runes := []rune(s)
	n := len(runes)
	if n <= 8 {
		return strings.Repeat("*", n)
	}
	return string(runes[:4]) + "…" + string(runes[n-4:])
}

func Boolyn(b bool) string {
	if b {
		return "✓"
	}
	return "✗"
}

func SendArrow(isSend int) string {
	if isSend == 1 {
		return "↑"
	}
	return "↓"
}

// IndentJSON 重新格式化任意 JSON（保留原字段顺序）。
func IndentJSON(b []byte) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, b, "", "  "); err != nil {
		return "", fmt.Errorf("格式化 JSON: %w", err)
	}
	return buf.String(), nil
}

type appShareXML struct {
	XMLName xml.Name `xml:"msg"`
	AppMsg  struct {
		Title string `xml:"title"`
		Des   string `xml:"des"`
		URL   string `xml:"url"`
	} `xml:"appmsg"`
	AppInfo struct {
		AppName string `xml:"appname"`
	} `xml:"appinfo"`
}

type revokeXML struct {
	XMLName  xml.Name `xml:"sysmsg"`
	RevokeMsg struct {
		Content string `xml:"content"`
	} `xml:"revokemsg"`
}

// FormatContent 把微信 XML 卡片（链接/小程序/B站分享等）和撤回消息提取成简洁文本。
// 非 XML 内容原样返回。
func FormatContent(content string) string {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "<?xml") && !strings.HasPrefix(trimmed, "<msg>") {
		return content
	}
	// 撤回消息
	var rv revokeXML
	if err := xml.Unmarshal([]byte(trimmed), &rv); err == nil && rv.RevokeMsg.Content != "" {
		return "[" + strings.TrimSpace(rv.RevokeMsg.Content) + "]"
	}
	// 分享卡片
	var m appShareXML
	if err := xml.Unmarshal([]byte(trimmed), &m); err != nil {
		return content
	}
	title := strings.TrimSpace(m.AppMsg.Title)
	url := strings.TrimSpace(m.AppMsg.URL)
	appName := strings.TrimSpace(m.AppInfo.AppName)
	if title == "" && url == "" {
		return content
	}
	var b strings.Builder
	if appName != "" {
		b.WriteString("【")
		b.WriteString(appName)
		b.WriteString("】")
	}
	if title != "" {
		b.WriteString(title)
	}
	if url != "" {
		if b.Len() > 0 {
			b.WriteString(" ")
		}
		b.WriteString(url)
	}
	return b.String()
}

// ConfigSummary 把 Config 转成 (key, value) 对列表，供 table 渲染。
func ConfigSummary(c *config.Config) [][2]string {
	src := c.SourcePath
	if src == "" {
		src = "(无；使用默认值/环境变量/flag)"
	}
	return [][2]string{
		{"source", src},
		{"host", c.Host},
		{"token", MaskToken(c.Token)},
		{"timeout", c.Timeout.String()},
		{"output.color", Boolyn(c.Output.Color)},
	}
}
