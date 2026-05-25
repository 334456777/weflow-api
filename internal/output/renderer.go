package output

import (
	"encoding/json"

	"github.com/334456777/weflow-api/internal/client"
	"github.com/334456777/weflow-api/internal/config"
)

type Renderer interface {
	Health(r *client.HealthResponse) error

	Sessions(r *client.SessionsResponse) error
	SessionsChatLab(r *client.ChatLabSessionsResponse) error

	Messages(r *client.MessagesResponse) error
	MessagesChatLab(r *client.ChatLabMessagesResponse) error

	Contacts(r *client.ContactsResponse) error
	Members(r *client.GroupMembersResponse) error

	Event(e client.Event) error
	Config(c *config.Config) error

	Raw(b json.RawMessage) error
}
