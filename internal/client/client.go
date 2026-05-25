package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/334456777/weflow-api/internal/config"
)

const userAgent = "weflow-cli"

type Client struct {
	baseURL    *url.URL
	token      string
	httpClient *http.Client
}

func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("配置中 host 为空，请运行 `weflow config init` 或设置 --host / WEFLOW_HOST")
	}
	u, err := url.Parse(cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("解析 host %q: %w", cfg.Host, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("host %q 不是合法 URL（需要 scheme + host）", cfg.Host)
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = config.DefaultTimeout
	}
	return &Client{
		baseURL:    u,
		token:      cfg.Token,
		httpClient: &http.Client{Timeout: timeout},
	}, nil
}

func (c *Client) BaseURL() string { return c.baseURL.String() }
func (c *Client) Token() string   { return c.token }

// SetHTTPClient 允许子命令（如 sns export 需要更长超时）替换底层 client。
func (c *Client) SetHTTPClient(h *http.Client) { c.httpClient = h }

// DoJSON 发起 JSON 请求并解码响应到 out（nil 表示丢弃响应体）。
func (c *Client) DoJSON(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	req, err := c.newRequest(ctx, method, path, query, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求 %s 失败: %w", path, err)
	}
	defer resp.Body.Close()
	return decodeResponse(resp, out)
}

// DoStream 发起请求并把响应交给调用方（调用方必须 Close Body）。用于 SSE 与媒体下载。
func (c *Client) DoStream(ctx context.Context, method, path string, query url.Values) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, query, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 %s 失败: %w", path, err)
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, decodeAPIError(resp)
	}
	return resp, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body any) (*http.Request, error) {
	u := *c.baseURL
	u.Path = strings.TrimRight(u.Path, "/") + path
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("构造请求: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return req, nil
}

// BuildURL 拼出绝对 URL（用于打印或第三方下载工具）。
func (c *Client) BuildURL(path string, query url.Values) string {
	u := *c.baseURL
	u.Path = strings.TrimRight(u.Path, "/") + path
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String()
}

func decodeResponse(resp *http.Response, out any) error {
	if resp.StatusCode >= 400 {
		return decodeAPIError(resp)
	}
	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("解析响应: %w", err)
	}
	return nil
}
