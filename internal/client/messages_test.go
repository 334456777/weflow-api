package client

import (
	"strconv"
	"testing"
	"time"
)

func TestExpandEndDay(t *testing.T) {
	got := expandEndDay("20260527")
	ts, err := strconv.ParseInt(got, 10, 64)
	if err != nil {
		t.Fatalf("expected unix timestamp, got %q: %v", got, err)
	}
	want := time.Date(2026, 5, 27, 23, 59, 59, 0, time.Local).Unix()
	if ts != want {
		t.Errorf("expandEndDay(20260527) = %d, want %d", ts, want)
	}

	if got := expandEndDay("1738713600"); got != "1738713600" {
		t.Errorf("plain timestamp should pass through, got %q", got)
	}
	if got := expandEndDay(""); got != "" {
		t.Errorf("empty should pass through, got %q", got)
	}
	if got := expandEndDay("2026-05-27"); got != "2026-05-27" {
		t.Errorf("non-8-digit string should pass through, got %q", got)
	}
}
