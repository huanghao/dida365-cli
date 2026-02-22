package cli

import (
	"fmt"
	"strings"
)

const (
	outputFormatTable = "table"
	outputFormatJSON  = "json"
)

func resolveOutputFormat(format string, asJSON bool) (string, error) {
	f := strings.ToLower(strings.TrimSpace(format))
	if asJSON {
		if f != "" && f != outputFormatJSON {
			return "", fmt.Errorf("--json cannot be combined with --format=%s", format)
		}
		return outputFormatJSON, nil
	}
	if f == "" {
		return outputFormatTable, nil
	}
	switch f {
	case outputFormatTable, outputFormatJSON:
		return f, nil
	default:
		return "", fmt.Errorf("unsupported format %q (use table or json)", format)
	}
}
