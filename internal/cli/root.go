package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRoot(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dida",
		Short:         "Dida365 CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if app == nil {
				return fmt.Errorf("app not initialized")
			}
			return app.ReloadConfigStore()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().BoolVar(&app.Debug, "debug", false, "Enable debug output")
	cmd.PersistentFlags().BoolVar(&app.DryRun, "dry-run", false, "Print request intent without executing")
	cmd.PersistentFlags().StringVar(&app.ConfigPath, "config", "", "Path to config file (default: ~/.config/dida365-cli/config.json)")

	cmd.AddCommand(NewAuthCommand(app))
	cmd.AddCommand(NewProjectsCommand(app))
	cmd.AddCommand(NewListCommand(app))
	cmd.AddCommand(NewShowCommand(app))
	cmd.AddCommand(NewAddCommand(app))
	cmd.AddCommand(NewUpdateCommand(app))
	cmd.AddCommand(NewDoneCommand(app))
	cmd.AddCommand(NewDeleteCommand(app))
	cmd.AddCommand(NewVersionCommand(app))

	return cmd
}

func Execute() error {
	app, err := NewApp()
	if err != nil {
		return err
	}
	if err := NewRoot(app).Execute(); err != nil {
		return err
	}
	return nil
}
