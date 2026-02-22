package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func PrintJSON(w io.Writer, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

func PrintSimpleTable(w io.Writer, headers []string, rows [][]string) error {
	if len(headers) == 0 {
		return nil
	}
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	printRow := func(cols []string) error {
		parts := make([]string, 0, len(headers))
		for i := range headers {
			val := ""
			if i < len(cols) {
				val = cols[i]
			}
			parts = append(parts, padRight(val, widths[i]))
		}
		_, err := fmt.Fprintln(w, strings.Join(parts, "  "))
		return err
	}

	if err := printRow(headers); err != nil {
		return err
	}
	sep := make([]string, len(headers))
	for i := range headers {
		sep[i] = strings.Repeat("-", widths[i])
	}
	if err := printRow(sep); err != nil {
		return err
	}
	for _, row := range rows {
		if err := printRow(row); err != nil {
			return err
		}
	}
	return nil
}

func padRight(v string, width int) string {
	if len(v) >= width {
		return v
	}
	return v + strings.Repeat(" ", width-len(v))
}
