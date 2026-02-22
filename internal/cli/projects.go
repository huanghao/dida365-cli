package cli

import (
	"fmt"

	"github.com/huanghao/dida365-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewProjectsCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "Project operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newProjectsListCommand(app))
	return cmd
}

func newProjectsListCommand(app *App) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			if app.DryRun {
				fmt.Fprintf(app.Out, "Would call GET %s/project\n", cfg.APIBaseURL)
				return nil
			}
			client := newAPIClient(cfg)
			projects, err := client.GetProjects()
			if err != nil {
				return err
			}
			if asJSON {
				return output.PrintJSON(app.Out, projects)
			}
			rows := make([][]string, 0, len(projects))
			for _, p := range projects {
				rows = append(rows, []string{p.ID, p.Name, p.Kind, fmt.Sprintf("%t", p.Closed)})
			}
			return output.PrintSimpleTable(app.Out, []string{"ID", "Name", "Kind", "Closed"}, rows)
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}
