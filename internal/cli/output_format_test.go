package cli

import "testing"

func TestResolveOutputFormat_DefaultTable(t *testing.T) {
	got, err := resolveOutputFormat("", false)
	if err != nil {
		t.Fatalf("resolveOutputFormat returned error: %v", err)
	}
	if got != outputFormatTable {
		t.Fatalf("expected %q, got %q", outputFormatTable, got)
	}
}

func TestResolveOutputFormat_JSONFlag(t *testing.T) {
	got, err := resolveOutputFormat("", true)
	if err != nil {
		t.Fatalf("resolveOutputFormat returned error: %v", err)
	}
	if got != outputFormatJSON {
		t.Fatalf("expected %q, got %q", outputFormatJSON, got)
	}
}

func TestResolveOutputFormat_Conflict(t *testing.T) {
	if _, err := resolveOutputFormat("table", true); err == nil {
		t.Fatalf("expected conflict error")
	}
}

func TestResolveOutputFormat_Unsupported(t *testing.T) {
	if _, err := resolveOutputFormat("csv", false); err == nil {
		t.Fatalf("expected unsupported format error")
	}
}

func TestResolveOutputFormat_ExplicitJSON(t *testing.T) {
	got, err := resolveOutputFormat("json", false)
	if err != nil {
		t.Fatalf("resolveOutputFormat returned error: %v", err)
	}
	if got != outputFormatJSON {
		t.Fatalf("expected %q, got %q", outputFormatJSON, got)
	}
}
