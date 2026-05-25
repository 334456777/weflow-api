package output

import (
	"bytes"
	"encoding/json"
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
