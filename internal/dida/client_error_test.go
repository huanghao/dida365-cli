package dida

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetProjects_Unauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":401,"message":"Unauthorized"}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "bad-token")
	_, err := c.GetProjects()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "status=401") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":404,"message":"Not Found"}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "token")
	_, err := c.GetTask("p1", "t1")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "status=404") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetProjects_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "token")
	_, err := c.GetProjects()
	if err == nil {
		t.Fatalf("expected decode error")
	}
	if !strings.Contains(err.Error(), "decode response") {
		t.Fatalf("unexpected error: %v", err)
	}
}
