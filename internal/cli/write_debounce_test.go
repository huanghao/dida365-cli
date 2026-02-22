package cli

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/huanghao/dida365-cli/internal/config"
)

func TestWriteDebounceBlocksRapidDuplicate(t *testing.T) {
	tmp := t.TempDir()
	store, err := config.NewStore(filepath.Join(tmp, "config.json"))
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	app := &App{ConfigStore: store}
	payload := map[string]string{"project": "p1", "title": "hello"}

	if err := checkWriteDebounce(app, "add_task", payload); err != nil {
		t.Fatalf("unexpected pre-check error: %v", err)
	}

	markWriteDebounce(app, "add_task", payload)

	if err := checkWriteDebounce(app, "add_task", payload); err == nil {
		t.Fatalf("expected debounce error, got nil")
	}
}

func TestWriteDebounceAllowsAfterWindow(t *testing.T) {
	tmp := t.TempDir()
	store, err := config.NewStore(filepath.Join(tmp, "config.json"))
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	app := &App{ConfigStore: store}
	payload := map[string]string{"project": "p1", "title": "hello"}

	markWriteDebounce(app, "add_task", payload)

	path := debounceFilePath(app)
	state, err := loadDebounceState(path)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	key, err := debounceKey("add_task", payload)
	if err != nil {
		t.Fatalf("debounce key: %v", err)
	}
	state.Entries[key] = time.Now().Add(-writeDebounceWindow - time.Second).Unix()
	if err := saveDebounceState(path, state); err != nil {
		t.Fatalf("save state: %v", err)
	}

	if err := checkWriteDebounce(app, "add_task", payload); err != nil {
		t.Fatalf("unexpected error after window: %v", err)
	}
}
