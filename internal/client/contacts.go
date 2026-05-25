package client

import (
	"context"
	"net/url"
	"strconv"
)

type Contact struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Remark      string `json:"remark"`
	Nickname    string `json:"nickname"`
	Alias       string `json:"alias"`
	AvatarURL   string `json:"avatarUrl"`
	Type        string `json:"type"`
}

type ContactsResponse struct {
	Success  bool      `json:"success"`
	Count    int       `json:"count"`
	Contacts []Contact `json:"contacts"`
}

type ContactsParams struct {
	Keyword string
	Limit   int
}

func (p ContactsParams) toQuery() url.Values {
	q := url.Values{}
	if p.Keyword != "" {
		q.Set("keyword", p.Keyword)
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	return q
}

func (c *Client) Contacts(ctx context.Context, p ContactsParams) (*ContactsResponse, error) {
	var out ContactsResponse
	if err := c.DoJSON(ctx, "GET", "/api/v1/contacts", p.toQuery(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
