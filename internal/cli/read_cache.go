package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const readCacheTTL = 10 * time.Second

type readCacheState struct {
	Entries map[string]readCacheEntry `json:"entries"`
}

type readCacheEntry struct {
	ExpiresAt int64           `json:"expires_at"`
	Payload   json.RawMessage `json:"payload"`
}

func readCacheEnabled(app *App) bool {
	if app != nil && app.NoCache {
		return false
	}
	return strings.TrimSpace(os.Getenv("DIDA_NO_CACHE")) != "1"
}

func readCacheFilePath(app *App) string {
	if app == nil || app.ConfigStore == nil {
		return ""
	}
	cfgPath := app.ConfigStore.Path()
	if cfgPath == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(cfgPath), "cache.json")
}

func readCacheKey(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(sum[:])
}

func readCacheGet(app *App, key string, dst any) bool {
	if !readCacheEnabled(app) {
		return false
	}
	path := readCacheFilePath(app)
	if path == "" {
		return false
	}
	state, err := loadReadCacheState(path)
	if err != nil {
		return false
	}
	entry, ok := state.Entries[key]
	if !ok {
		return false
	}
	if time.Now().Unix() >= entry.ExpiresAt {
		delete(state.Entries, key)
		_ = saveReadCacheState(path, state)
		return false
	}
	if err := json.Unmarshal(entry.Payload, dst); err != nil {
		return false
	}
	return true
}

func readCachePut(app *App, key string, v any) {
	if !readCacheEnabled(app) {
		return
	}
	path := readCacheFilePath(app)
	if path == "" {
		return
	}
	payload, err := json.Marshal(v)
	if err != nil {
		return
	}
	state, err := loadReadCacheState(path)
	if err != nil {
		state = readCacheState{Entries: map[string]readCacheEntry{}}
	}
	if state.Entries == nil {
		state.Entries = map[string]readCacheEntry{}
	}

	now := time.Now().Unix()
	state.Entries[key] = readCacheEntry{
		ExpiresAt: now + int64(readCacheTTL.Seconds()),
		Payload:   payload,
	}
	for k, entry := range state.Entries {
		if now >= entry.ExpiresAt {
			delete(state.Entries, k)
		}
	}

	_ = saveReadCacheState(path, state)
}

func clearReadCache(app *App) {
	path := readCacheFilePath(app)
	if path == "" {
		return
	}
	_ = saveReadCacheState(path, readCacheState{Entries: map[string]readCacheEntry{}})
}

func loadReadCacheState(path string) (readCacheState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return readCacheState{Entries: map[string]readCacheEntry{}}, nil
		}
		return readCacheState{}, err
	}
	var state readCacheState
	if err := json.Unmarshal(data, &state); err != nil {
		return readCacheState{}, err
	}
	if state.Entries == nil {
		state.Entries = map[string]readCacheEntry{}
	}
	return state, nil
}

func saveReadCacheState(path string, state readCacheState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
}
