package cli

import (
	"fmt"

	"github.com/huanghao/dida365-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewListCommand(app *App) *cobra.Command {
	var projectID string
	var format string
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
			resolvedFormat, err := resolveOutputFormat(format, asJSON)
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
			if resolvedFormat == outputFormatJSON {
				return output.PrintJSON(app.Out, data)
			}
			rows := make([][]string, 0, len(data.Tasks))
			for _, t := range data.Tasks {
				content := ellipsis(t.Content, 40)
				rows = append(rows, []string{
					t.ID,
					t.Title,
					formatStatus(t.Status, t.CompletedTime),
					completeLabel(t.Status, t.CompletedTime),
					t.DueDate,
					formatPriority(t.Priority),
					content,
				})
			}
			return output.PrintSimpleTable(app.Out, []string{"ID", "Title", "Status", "Completed", "Due", "Priority", "Content"}, rows)
		},
	}

	cmd.Flags().StringVar(&projectID, "project", "", "Project ID")
	cmd.Flags().StringVar(&format, "format", "", "Output format: table, json")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}
