package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error,omitempty"`
	Raw        string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API 错误 %d: %s", e.StatusCode, e.Message)
	}
	if e.Raw != "" {
		return fmt.Sprintf("API 错误 %d: %s", e.StatusCode, strings.TrimSpace(e.Raw))
	}
	return fmt.Sprintf("API 错误 %d", e.StatusCode)
}

func decodeAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	e := &APIError{StatusCode: resp.StatusCode, Raw: string(body)}
	if len(body) > 0 {
		_ = json.Unmarshal(body, e)
	}
	return e
}
