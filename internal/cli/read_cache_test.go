package cli

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/huanghao/dida365-cli/internal/config"
)

func TestReadCachePutGetAndClear(t *testing.T) {
	tmp := t.TempDir()
	store, err := config.NewStore(filepath.Join(tmp, "config.json"))
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	app := &App{ConfigStore: store}
	key := readCacheKey("GET", "/project")
	in := map[string]string{"id": "p1", "name": "demo"}

	readCachePut(app, key, in)

	var out map[string]string
	if ok := readCacheGet(app, key, &out); !ok {
		t.Fatalf("expected cache hit")
	}
	if out["id"] != "p1" || out["name"] != "demo" {
		t.Fatalf("unexpected payload: %#v", out)
	}

	clearReadCache(app)
	if ok := readCacheGet(app, key, &out); ok {
		t.Fatalf("expected cache miss after clear")
	}
}

func TestReadCacheExpires(t *testing.T) {
	tmp := t.TempDir()
	store, err := config.NewStore(filepath.Join(tmp, "config.json"))
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	app := &App{ConfigStore: store}
	key := readCacheKey("GET", "/project", "p1", "data")
	readCachePut(app, key, map[string]string{"v": "x"})

	path := readCacheFilePath(app)
	state, err := loadReadCacheState(path)
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	entry := state.Entries[key]
	entry.ExpiresAt = time.Now().Add(-time.Second).Unix()
	state.Entries[key] = entry
	if err := saveReadCacheState(path, state); err != nil {
		t.Fatalf("save state: %v", err)
	}

	var out map[string]string
	if ok := readCacheGet(app, key, &out); ok {
		t.Fatalf("expected expired cache miss")
	}
}

func TestReadCacheDisabledByEnvAndFlag(t *testing.T) {
	tmp := t.TempDir()
	store, err := config.NewStore(filepath.Join(tmp, "config.json"))
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	app := &App{ConfigStore: store}
	key := readCacheKey("GET", "/project")

	t.Setenv("DIDA_NO_CACHE", "1")
	readCachePut(app, key, map[string]string{"v": "x"})
	var out map[string]string
	if ok := readCacheGet(app, key, &out); ok {
		t.Fatalf("expected miss when DIDA_NO_CACHE=1")
	}

	t.Setenv("DIDA_NO_CACHE", "")
	app.NoCache = true
	readCachePut(app, key, map[string]string{"v": "y"})
	if ok := readCacheGet(app, key, &out); ok {
		t.Fatalf("expected miss when --no-cache enabled")
	}
}
