package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/334456777/weflow-api/internal/config"
)

func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := New(&config.Config{Host: srv.URL, Token: "test-token", Timeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestDoJSON_Success(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q", got)
		}
		if r.URL.Path != "/api/v1/health" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	}))

	var out struct {
		Status string `json:"status"`
	}
	if err := c.DoJSON(context.Background(), "GET", "/api/v1/health", nil, nil, &out); err != nil {
		t.Fatalf("DoJSON: %v", err)
	}
	if out.Status != "ok" {
		t.Fatalf("status = %q", out.Status)
	}
}

func TestDoJSON_APIError(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error":"Media not found"}`)
	}))

	err := c.DoJSON(context.Background(), "GET", "/api/v1/media/foo", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("status = %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Media not found" {
		t.Errorf("message = %q", apiErr.Message)
	}
}

func TestDoJSON_RawErrorBody(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `boom`)
	}))

	err := c.DoJSON(context.Background(), "GET", "/", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") || !strings.Contains(err.Error(), "boom") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDoJSON_PostBody(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q", ct)
		}
		buf := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(buf)
		if !strings.Contains(string(buf), `"foo":"bar"`) {
			t.Errorf("body = %s", buf)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"ok":true}`)
	}))

	body := map[string]string{"foo": "bar"}
	if err := c.DoJSON(context.Background(), "POST", "/api/v1/x", nil, body, nil); err != nil {
		t.Fatalf("DoJSON: %v", err)
	}
}

func TestBuildURL(t *testing.T) {
	c, err := New(&config.Config{Host: "http://127.0.0.1:5031", Timeout: time.Second})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	got := c.BuildURL("/api/v1/push/messages", nil)
	if got != "http://127.0.0.1:5031/api/v1/push/messages" {
		t.Errorf("BuildURL = %q", got)
	}
}

func TestStreamEvents(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		fmt.Fprint(w, "event: message.new\ndata: {\"a\":1}\n\n")
		flusher.Flush()
		fmt.Fprint(w, "event: message.revoke\ndata: {\"b\":2}\n\n")
		flusher.Flush()
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	events, errs := c.StreamEvents(ctx, "/api/v1/push/messages", nil)

	var got []Event
	for e := range events {
		got = append(got, e)
	}
	select {
	case err := <-errs:
		if err != nil {
			t.Fatalf("stream err: %v", err)
		}
	default:
	}
	if len(got) != 2 {
		t.Fatalf("got %d events, want 2: %+v", len(got), got)
	}
	if got[0].Name != "message.new" || got[0].Data != `{"a":1}` {
		t.Errorf("event[0] = %+v", got[0])
	}
	if got[1].Name != "message.revoke" || got[1].Data != `{"b":2}` {
		t.Errorf("event[1] = %+v", got[1])
	}
}

func TestNew_InvalidHost(t *testing.T) {
	cases := []struct {
		name string
		host string
	}{
		{"empty", ""},
		{"no scheme", "127.0.0.1:5031"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := New(&config.Config{Host: tc.host}); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
