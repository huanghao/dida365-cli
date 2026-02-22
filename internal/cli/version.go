package cli

import (
	"fmt"

	"github.com/huanghao/dida365-cli/internal/output"
	"github.com/spf13/cobra"
)

var Version = "dev"

func NewVersionCommand(app *App) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if asJSON {
				return output.PrintJSON(app.Out, map[string]any{
					"version": Version,
				})
			}
			fmt.Fprintln(app.Out, Version)
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}
