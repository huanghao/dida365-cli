package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

func NewVersionCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(app.Out, Version)
			return nil
		},
	}
}
