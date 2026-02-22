package cli

import "strings"

func taskCompleted(status int, completedTime string) bool {
	if strings.TrimSpace(completedTime) != "" {
		return true
	}
	// Keep a pragmatic fallback for status-only responses.
	return status != 0
}

func completeLabel(status int, completedTime string) string {
	if taskCompleted(status, completedTime) {
		return "yes"
	}
	return "no"
}

func ellipsis(v string, max int) string {
	if max <= 0 {
		return v
	}
	runes := []rune(v)
	if len(runes) <= max {
		return v
	}
	if max <= 1 {
		return string(runes[:max])
	}
	return string(runes[:max-1]) + "…"
}
