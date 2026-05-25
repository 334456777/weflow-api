package client

import (
	"context"
	"net/url"
)

type GroupMember struct {
	WxID          string `json:"wxid"`
	DisplayName   string `json:"displayName"`
	Nickname      string `json:"nickname"`
	Remark        string `json:"remark"`
	Alias         string `json:"alias"`
	GroupNickname string `json:"groupNickname"`
	AvatarURL     string `json:"avatarUrl"`
	IsOwner       bool   `json:"isOwner"`
	IsFriend      bool   `json:"isFriend"`
	MessageCount  int    `json:"messageCount,omitempty"`
}

type GroupMembersResponse struct {
	Success    bool          `json:"success"`
	ChatroomID string        `json:"chatroomId"`
	Count      int           `json:"count"`
	FromCache  bool          `json:"fromCache"`
	UpdatedAt  int64         `json:"updatedAt"`
	Members    []GroupMember `json:"members"`
}

type MembersParams struct {
	ChatroomID   string
	WithCounts   bool
	ForceRefresh bool
}

func (p MembersParams) toQuery() url.Values {
	q := url.Values{}
	q.Set("chatroomId", p.ChatroomID)
	if p.WithCounts {
		q.Set("includeMessageCounts", "1")
	}
	if p.ForceRefresh {
		q.Set("forceRefresh", "1")
	}
	return q
}

func (c *Client) Members(ctx context.Context, p MembersParams) (*GroupMembersResponse, error) {
	var out GroupMembersResponse
	if err := c.DoJSON(ctx, "GET", "/api/v1/group-members", p.toQuery(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
