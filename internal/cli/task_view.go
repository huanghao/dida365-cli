package cli

import "strings"

func taskCompleted(status int, completedTime string) bool {
	if strings.TrimSpace(completedTime) != "" {
		return true
	}
	// Dida docs/examples use status=0 for incomplete; treat 2 as completed fallback.
	return status == 2
}

func completeLabel(status int, completedTime string) string {
	if taskCompleted(status, completedTime) {
		return "yes"
	}
	return "no"
}

func statusLabel(status int, completedTime string) string {
	if taskCompleted(status, completedTime) {
		return "completed"
	}
	switch status {
	case 0:
		return "incomplete"
	default:
		return "status_unknown"
	}
}

func formatStatus(status int, completedTime string) string {
	return statusLabel(status, completedTime)
}

func priorityLabel(priority int) string {
	switch priority {
	case 0:
		return "none"
	case 1:
		return "low"
	case 3:
		return "medium"
	case 5:
		return "high"
	default:
		return "priority_unknown"
	}
}

func formatPriority(priority int) string {
	return priorityLabel(priority)
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
