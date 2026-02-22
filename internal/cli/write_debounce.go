package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const writeDebounceWindow = 3 * time.Second

type debounceState struct {
	Entries map[string]int64 `json:"entries"`
}

func checkWriteDebounce(app *App, action string, payload any) error {
	path := debounceFilePath(app)
	if path == "" {
		return nil
	}
	key, err := debounceKey(action, payload)
	if err != nil {
		return err
	}
	state, err := loadDebounceState(path)
	if err != nil {
		// fail-open: debounce is a safety optimization and should not block writes
		return nil
	}
	if ts, ok := state.Entries[key]; ok {
		if time.Since(time.Unix(ts, 0)) < writeDebounceWindow {
			return fmt.Errorf("blocked duplicate write within %ds; wait and retry", int(writeDebounceWindow.Seconds()))
		}
	}
	return nil
}

func markWriteDebounce(app *App, action string, payload any) {
	path := debounceFilePath(app)
	if path == "" {
		return
	}
	key, err := debounceKey(action, payload)
	if err != nil {
		return
	}

	state, err := loadDebounceState(path)
	if err != nil {
		state = debounceState{Entries: map[string]int64{}}
	}
	if state.Entries == nil {
		state.Entries = map[string]int64{}
	}

	now := time.Now().Unix()
	state.Entries[key] = now

	// Keep file small: remove stale keys older than one hour.
	cutoff := now - 3600
	for k, ts := range state.Entries {
		if ts < cutoff {
			delete(state.Entries, k)
		}
	}

	_ = saveDebounceState(path, state)
}

func debounceFilePath(app *App) string {
	if app == nil || app.ConfigStore == nil {
		return ""
	}
	cfgPath := app.ConfigStore.Path()
	if cfgPath == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(cfgPath), "debounce.json")
}

func debounceKey(action string, payload any) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(append([]byte(action+"|"), data...))
	return hex.EncodeToString(sum[:]), nil
}

func loadDebounceState(path string) (debounceState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return debounceState{Entries: map[string]int64{}}, nil
		}
		return debounceState{}, err
	}
	var state debounceState
	if err := json.Unmarshal(data, &state); err != nil {
		return debounceState{}, err
	}
	if state.Entries == nil {
		state.Entries = map[string]int64{}
	}
	return state, nil
}

func saveDebounceState(path string, state debounceState) error {
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
