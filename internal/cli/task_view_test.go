package cli

import "testing"

func TestCompleteLabel(t *testing.T) {
	if got := completeLabel(0, ""); got != "no" {
		t.Fatalf("expected no, got %s", got)
	}
	if got := completeLabel(0, "2026-02-22T00:00:00+0000"); got != "yes" {
		t.Fatalf("expected yes, got %s", got)
	}
	if got := completeLabel(2, ""); got != "yes" {
		t.Fatalf("expected yes, got %s", got)
	}
}

func TestEllipsis(t *testing.T) {
	if got := ellipsis("short", 10); got != "short" {
		t.Fatalf("unexpected result: %s", got)
	}
	if got := ellipsis("1234567890", 5); got != "1234…" {
		t.Fatalf("unexpected result: %s", got)
	}
}
