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

func TestStatusLabel(t *testing.T) {
	if got := formatStatus(0, ""); got != "incomplete" {
		t.Fatalf("expected incomplete, got %s", got)
	}
	if got := formatStatus(2, ""); got != "completed" {
		t.Fatalf("expected completed, got %s", got)
	}
	if got := formatStatus(9, ""); got != "status_unknown" {
		t.Fatalf("expected status_unknown, got %s", got)
	}
}

func TestPriorityLabel(t *testing.T) {
	cases := map[int]string{
		0: "none",
		1: "low",
		3: "medium",
		5: "high",
		9: "priority_unknown",
	}
	for in, expected := range cases {
		if got := formatPriority(in); got != expected {
			t.Fatalf("priority %d expected %s, got %s", in, expected, got)
		}
	}
}
