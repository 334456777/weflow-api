package client

import (
	"context"
	"net/url"
	"strconv"
)

type PullParams struct {
	Since  int64
	End    int64
	Limit  int
	Offset int
}

func (p PullParams) toQuery() url.Values {
	q := url.Values{}
	if p.Since > 0 {
		q.Set("since", strconv.FormatInt(p.Since, 10))
	}
	if p.End > 0 {
		q.Set("end", strconv.FormatInt(p.End, 10))
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset > 0 {
		q.Set("offset", strconv.Itoa(p.Offset))
	}
	return q
}

// Pull 调用 §4.2 GET /api/v1/sessions/:id/messages
func (c *Client) Pull(ctx context.Context, sessionID string, p PullParams) (*ChatLabMessagesResponse, error) {
	var out ChatLabMessagesResponse
	path := "/api/v1/sessions/" + url.PathEscape(sessionID) + "/messages"
	if err := c.DoJSON(ctx, "GET", path, p.toQuery(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
