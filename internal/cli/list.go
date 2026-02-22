package cli

import (
	"fmt"

	"github.com/huanghao/dida365-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewListCommand(app *App) *cobra.Command {
	var projectID string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" {
				return fmt.Errorf("--project is required")
			}
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			if app.DryRun {
				fmt.Fprintf(app.Out, "Would call GET %s/project/%s/data\n", cfg.APIBaseURL, projectID)
				return nil
			}
			client := newAPIClient(cfg)
			data, err := client.GetProjectData(projectID)
			if err != nil {
				return err
			}
			if asJSON {
				return output.PrintJSON(app.Out, data)
			}
			rows := make([][]string, 0, len(data.Tasks))
			for _, t := range data.Tasks {
				rows = append(rows, []string{t.ID, t.Title, fmt.Sprintf("%d", t.Status), t.DueDate})
			}
			return output.PrintSimpleTable(app.Out, []string{"ID", "Title", "Status", "Due"}, rows)
		},
	}

	cmd.Flags().StringVar(&projectID, "project", "", "Project ID")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}
