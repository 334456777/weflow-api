package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/334456777/weflow-api/internal/client"
	"github.com/334456777/weflow-api/internal/config"
)

type JSONRenderer struct {
	w io.Writer
}

func NewJSON(w io.Writer) *JSONRenderer { return &JSONRenderer{w: w} }

func (r *JSONRenderer) write(v any) error {
	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("编码 JSON: %w", err)
	}
	return nil
}

func (r *JSONRenderer) Health(v *client.HealthResponse) error              { return r.write(v) }
func (r *JSONRenderer) Sessions(v *client.SessionsResponse) error          { return r.write(v) }
func (r *JSONRenderer) SessionsChatLab(v *client.ChatLabSessionsResponse) error {
	return r.write(v)
}
func (r *JSONRenderer) Messages(v *client.MessagesResponse) error { return r.write(v) }
func (r *JSONRenderer) MessagesChatLab(v *client.ChatLabMessagesResponse) error {
	return r.write(v)
}
func (r *JSONRenderer) Contacts(v *client.ContactsResponse) error    { return r.write(v) }
func (r *JSONRenderer) Members(v *client.GroupMembersResponse) error { return r.write(v) }
func (r *JSONRenderer) Event(e client.Event) error                   { return r.write(e) }
func (r *JSONRenderer) Config(c *config.Config) error                { return r.write(c) }

func (r *JSONRenderer) Raw(b json.RawMessage) error {
	s, err := IndentJSON(b)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(r.w, s)
	return err
}
