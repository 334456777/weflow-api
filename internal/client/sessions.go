package client

import (
	"context"
	"net/url"
	"strconv"
)

type Session struct {
	Username      string `json:"username"`
	DisplayName   string `json:"displayName"`
	Type          int    `json:"type"`
	LastTimestamp int64  `json:"lastTimestamp"`
	UnreadCount   int    `json:"unreadCount"`
}

type SessionsResponse struct {
	Success  bool      `json:"success"`
	Count    int       `json:"count"`
	Sessions []Session `json:"sessions"`
}

type ChatLabSession struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Platform      string `json:"platform"`
	Type          string `json:"type"`
	MessageCount  int    `json:"messageCount"`
	LastMessageAt int64  `json:"lastMessageAt"`
}

type ChatLabSessionsResponse struct {
	Sessions []ChatLabSession `json:"sessions"`
}

type SessionsParams struct {
	Keyword string
	Limit   int
}

func (p SessionsParams) toQuery(format string) url.Values {
	q := url.Values{}
	if p.Keyword != "" {
		q.Set("keyword", p.Keyword)
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if format != "" {
		q.Set("format", format)
	}
	return q
}

func (c *Client) Sessions(ctx context.Context, p SessionsParams) (*SessionsResponse, error) {
	var out SessionsResponse
	if err := c.DoJSON(ctx, "GET", "/api/v1/sessions", p.toQuery(""), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SessionsChatLab(ctx context.Context, p SessionsParams) (*ChatLabSessionsResponse, error) {
	var out ChatLabSessionsResponse
	if err := c.DoJSON(ctx, "GET", "/api/v1/sessions", p.toQuery("chatlab"), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
