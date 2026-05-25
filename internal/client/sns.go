package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

// === §7.1 时间线 ===
type SnsTimelineParams struct {
	Limit     int
	Offset    int
	Usernames string
	Keyword   string
	Start     string
	End       string
	Media     *int
	Replace   *int
	Inline    *int
}

func addOptInt(q url.Values, k string, p *int) {
	if p == nil {
		return
	}
	q.Set(k, strconv.Itoa(*p))
}

func (p SnsTimelineParams) toQuery() url.Values {
	q := url.Values{}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset > 0 {
		q.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Usernames != "" {
		q.Set("usernames", p.Usernames)
	}
	if p.Keyword != "" {
		q.Set("keyword", p.Keyword)
	}
	if p.Start != "" {
		q.Set("start", p.Start)
	}
	if p.End != "" {
		q.Set("end", p.End)
	}
	addOptInt(q, "media", p.Media)
	addOptInt(q, "replace", p.Replace)
	addOptInt(q, "inline", p.Inline)
	return q
}

func (c *Client) SnsTimeline(ctx context.Context, p SnsTimelineParams) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.DoJSON(ctx, "GET", "/api/v1/sns/timeline", p.toQuery(), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// === §7.2 发布者 ===
func (c *Client) SnsUsernames(ctx context.Context) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.DoJSON(ctx, "GET", "/api/v1/sns/usernames", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// === §7.3 导出统计 ===
type SnsExportStatsParams struct {
	Fast bool
}

func (c *Client) SnsExportStats(ctx context.Context, p SnsExportStatsParams) (json.RawMessage, error) {
	q := url.Values{}
	if p.Fast {
		q.Set("fast", "1")
	}
	var out json.RawMessage
	if err := c.DoJSON(ctx, "GET", "/api/v1/sns/export/stats", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// === §7.4 朋友圈媒体代理（流式下载） ===
func (c *Client) SnsProxy(ctx context.Context, rawURL, key string, dst io.Writer) (string, error) {
	q := url.Values{}
	q.Set("url", rawURL)
	if key != "" {
		q.Set("key", key)
	}
	resp, err := c.DoStream(ctx, "GET", "/api/v1/sns/media/proxy", q)
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

// === §7.5 导出朋友圈（同步 POST） ===
type SnsExportRequest struct {
	OutputDir        string `json:"outputDir"`
	Format           string `json:"format,omitempty"`
	Usernames        string `json:"usernames,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	ExportMedia      bool   `json:"exportMedia,omitempty"`
	ExportImages     bool   `json:"exportImages,omitempty"`
	ExportLivePhotos bool   `json:"exportLivePhotos,omitempty"`
	ExportVideos     bool   `json:"exportVideos,omitempty"`
	Start            string `json:"start,omitempty"`
	End              string `json:"end,omitempty"`
}

func (c *Client) SnsExport(ctx context.Context, req SnsExportRequest) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.DoJSON(ctx, "POST", "/api/v1/sns/export", nil, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// === §7.6 防删开关 ===
func (c *Client) SnsBlockDeleteStatus(ctx context.Context) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.DoJSON(ctx, "GET", "/api/v1/sns/block-delete/status", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) SnsBlockDeleteInstall(ctx context.Context) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.DoJSON(ctx, "POST", "/api/v1/sns/block-delete/install", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) SnsBlockDeleteUninstall(ctx context.Context) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.DoJSON(ctx, "POST", "/api/v1/sns/block-delete/uninstall", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// === §7.7 删除单条朋友圈 ===
func (c *Client) SnsDeletePost(ctx context.Context, postID string) (json.RawMessage, error) {
	var out json.RawMessage
	path := "/api/v1/sns/post/" + url.PathEscape(postID)
	if err := c.DoJSON(ctx, "DELETE", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
