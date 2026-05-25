package client

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Event 是一条 SSE 帧解析后的结构。Data 由多行 data: 用 \n 拼接而成。
type Event struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
	Data string `json:"data"`
}

// StreamEvents 订阅 SSE 长连接，按帧推送事件到 events 通道；连接断开或 ctx 取消时 events 关闭。
// 调用方应同时 select events 与 errs：errs 缓冲 1 条非 nil 错误后关闭。
func (c *Client) StreamEvents(ctx context.Context, path string, query url.Values) (<-chan Event, <-chan error) {
	events := make(chan Event)
	errs := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errs)

		resp, err := c.DoStream(ctx, "GET", path, query)
		if err != nil {
			errs <- err
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		var cur Event
		var dataLines []string
		flush := func() {
			if cur.Name == "" && cur.ID == "" && len(dataLines) == 0 {
				return
			}
			cur.Data = strings.Join(dataLines, "\n")
			select {
			case events <- cur:
			case <-ctx.Done():
			}
			cur = Event{}
			dataLines = nil
		}

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				flush()
				continue
			}
			if strings.HasPrefix(line, ":") {
				continue
			}
			field, value, ok := strings.Cut(line, ":")
			if !ok {
				field = line
				value = ""
			}
			value = strings.TrimPrefix(value, " ")
			switch field {
			case "event":
				cur.Name = value
			case "data":
				dataLines = append(dataLines, value)
			case "id":
				cur.ID = value
			case "retry":
				// 由上层 watch --reconnect 实现退避，此处忽略
			}
		}
		flush()

		if err := scanner.Err(); err != nil {
			select {
			case errs <- fmt.Errorf("读取 SSE 流: %w", err):
			case <-ctx.Done():
			}
		}
	}()

	return events, errs
}
