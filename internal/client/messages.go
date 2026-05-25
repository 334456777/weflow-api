package client

import (
	"context"
	"net/url"
	"strconv"
)

type Message struct {
	LocalID          int64         `json:"localId"`
	ServerID         string        `json:"serverId,omitempty"`
	LocalType        int64         `json:"localType"`
	CreateTime       int64         `json:"createTime"`
	IsSend           int           `json:"isSend"`
	SenderUsername   string        `json:"senderUsername"`
	Content          string        `json:"content"`
	RawContent       string        `json:"rawContent,omitempty"`
	ParsedContent    string        `json:"parsedContent,omitempty"`
	ReplyToMessageID string        `json:"replyToMessageId,omitempty"`
	Quote            *MessageQuote `json:"quote,omitempty"`
	MediaType        string        `json:"mediaType,omitempty"`
	MediaFileName    string        `json:"mediaFileName,omitempty"`
	MediaURL         string        `json:"mediaUrl,omitempty"`
	MediaLocalPath   string        `json:"mediaLocalPath,omitempty"`
}

type MessageQuote struct {
	PlatformMessageID string `json:"platformMessageId"`
	Sender            string `json:"sender"`
	AccountName       string `json:"accountName"`
	Content           string `json:"content"`
	Type              int    `json:"type"`
}

type MessagesMedia struct {
	Enabled    bool   `json:"enabled"`
	ExportPath string `json:"exportPath,omitempty"`
	Count      int    `json:"count"`
}

type MessagesResponse struct {
	Success  bool          `json:"success"`
	Talker   string        `json:"talker"`
	Count    int           `json:"count"`
	HasMore  bool          `json:"hasMore"`
	Media    MessagesMedia `json:"media"`
	Messages []Message     `json:"messages"`
}

type ChatLabInfo struct {
	Version    string `json:"version"`
	ExportedAt int64  `json:"exportedAt"`
	Generator  string `json:"generator"`
}

type ChatLabMeta struct {
	Name        string `json:"name"`
	Platform    string `json:"platform"`
	Type        string `json:"type"`
	GroupID     string `json:"groupId,omitempty"`
	GroupAvatar string `json:"groupAvatar,omitempty"`
	OwnerID     string `json:"ownerId,omitempty"`
}

type ChatLabMember struct {
	PlatformID    string `json:"platformId"`
	AccountName   string `json:"accountName"`
	GroupNickname string `json:"groupNickname,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
}

type ChatLabMessage struct {
	Sender            string `json:"sender"`
	AccountName       string `json:"accountName"`
	GroupNickname     string `json:"groupNickname,omitempty"`
	Timestamp         int64  `json:"timestamp"`
	Type              int    `json:"type"`
	Content           string `json:"content"`
	PlatformMessageID string `json:"platformMessageId,omitempty"`
	ReplyToMessageID  string `json:"replyToMessageId,omitempty"`
	MediaPath         string `json:"mediaPath,omitempty"`
}

type ChatLabSync struct {
	HasMore    bool  `json:"hasMore"`
	NextSince  int64 `json:"nextSince"`
	NextOffset int   `json:"nextOffset"`
	Watermark  int64 `json:"watermark"`
}

type ChatLabMessagesResponse struct {
	ChatLab  ChatLabInfo      `json:"chatlab"`
	Meta     ChatLabMeta      `json:"meta"`
	Members  []ChatLabMember  `json:"members"`
	Messages []ChatLabMessage `json:"messages"`
	Sync     *ChatLabSync     `json:"sync,omitempty"`
}

type MessagesParams struct {
	Talker  string
	Limit   int
	Offset  int
	Start   string
	End     string
	Keyword string
	Media   bool
	Image   *bool
	Voice   *bool
	Video   *bool
	Emoji   *bool
}

func addOptBool(q url.Values, k string, b *bool) {
	if b == nil {
		return
	}
	if *b {
		q.Set(k, "1")
	} else {
		q.Set(k, "0")
	}
}

func (p MessagesParams) toQuery(chatlab bool) url.Values {
	q := url.Values{}
	if p.Talker != "" {
		q.Set("talker", p.Talker)
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset > 0 {
		q.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Start != "" {
		q.Set("start", p.Start)
	}
	if p.End != "" {
		q.Set("end", p.End)
	}
	if p.Keyword != "" {
		q.Set("keyword", p.Keyword)
	}
	if chatlab {
		q.Set("chatlab", "1")
	}
	if p.Media {
		q.Set("media", "1")
		addOptBool(q, "image", p.Image)
		addOptBool(q, "voice", p.Voice)
		addOptBool(q, "video", p.Video)
		addOptBool(q, "emoji", p.Emoji)
	}
	return q
}

func (c *Client) Messages(ctx context.Context, p MessagesParams) (*MessagesResponse, error) {
	var out MessagesResponse
	if err := c.DoJSON(ctx, "GET", "/api/v1/messages", p.toQuery(false), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) MessagesChatLab(ctx context.Context, p MessagesParams) (*ChatLabMessagesResponse, error) {
	var out ChatLabMessagesResponse
	if err := c.DoJSON(ctx, "GET", "/api/v1/messages", p.toQuery(true), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
