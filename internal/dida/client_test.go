package dida

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildAuthorizeURL(t *testing.T) {
	u, err := BuildAuthorizeURL("cid", "http://localhost/callback", "tasks:read tasks:write", "s1")
	if err != nil {
		t.Fatalf("BuildAuthorizeURL returned error: %v", err)
	}
	if u == "" {
		t.Fatalf("expected non-empty URL")
	}
}

func TestGetProjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/project" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer token123" {
			t.Fatalf("missing bearer token")
		}
		_ = json.NewEncoder(w).Encode([]Project{{ID: "p1", Name: "Inbox"}})
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "token123")
	projects, err := c.GetProjects()
	if err != nil {
		t.Fatalf("GetProjects returned error: %v", err)
	}
	if len(projects) != 1 || projects[0].ID != "p1" {
		t.Fatalf("unexpected projects response: %+v", projects)
	}
}

func TestCreateTaskMissingToken(t *testing.T) {
	c := NewClient("https://api.dida365.com/open/v1", "")
	_, err := c.CreateTask(Task{ProjectID: "p1", Title: "t1"})
	if err == nil {
		t.Fatalf("expected error for missing token")
	}
}

func TestCreateProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/project" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if strings.TrimSpace(payload["name"].(string)) != "Demo Project" {
			t.Fatalf("unexpected name: %+v", payload)
		}
		_ = json.NewEncoder(w).Encode(Project{
			ID:       "p123",
			Name:     "Demo Project",
			Kind:     "TASK",
			ViewMode: "list",
			Color:    "#F18181",
		})
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "token123")
	project, err := c.CreateProject(CreateProjectInput{
		Name:     "Demo Project",
		ViewMode: "list",
		Kind:     "TASK",
		Color:    "#F18181",
	})
	if err != nil {
		t.Fatalf("CreateProject returned error: %v", err)
	}
	if project.ID != "p123" || project.Name != "Demo Project" {
		t.Fatalf("unexpected response: %+v", project)
	}
}
