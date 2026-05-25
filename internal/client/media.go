package client

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// MediaGet 流式拉取一个媒体文件并写入 dst，返回 Content-Type。
func (c *Client) MediaGet(ctx context.Context, relativePath string, dst io.Writer) (string, error) {
	relativePath = strings.TrimLeft(relativePath, "/")
	resp, err := c.DoStream(ctx, "GET", "/api/v1/media/"+relativePath, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ct := resp.Header.Get("Content-Type")
	if _, err := io.Copy(dst, resp.Body); err != nil {
		return ct, fmt.Errorf("写入媒体: %w", err)
	}
	return ct, nil
}

// MediaURL 返回带 access_token 的可访问 URL（仅用于打印或拷给其他工具）。
func (c *Client) MediaURL(relativePath string) string {
	relativePath = strings.TrimLeft(relativePath, "/")
	q := url.Values{}
	if c.token != "" {
		q.Set("access_token", c.token)
	}
	return c.BuildURL("/api/v1/media/"+relativePath, q)
}
