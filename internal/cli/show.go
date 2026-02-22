package cli

import (
	"fmt"

	"github.com/huanghao/dida365-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewShowCommand(app *App) *cobra.Command {
	var projectID string
	var taskID string
	var format string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show task details",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID == "" || taskID == "" {
				return fmt.Errorf("--project and --id are required")
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
				fmt.Fprintf(app.Out, "Would call GET %s/project/%s/task/%s\n", cfg.APIBaseURL, projectID, taskID)
				return nil
			}
			client := newAPIClient(cfg)
			task, err := client.GetTask(projectID, taskID)
			if err != nil {
				return err
			}
			if resolvedFormat == outputFormatJSON {
				return output.PrintJSON(app.Out, task)
			}
			rows := [][]string{
				{"ID", task.ID},
				{"Project", task.ProjectID},
				{"Title", task.Title},
				{"Status", fmt.Sprintf("%d", task.Status)},
				{"Start", task.StartDate},
				{"Due", task.DueDate},
				{"Priority", fmt.Sprintf("%d", task.Priority)},
			}
			return output.PrintSimpleTable(app.Out, []string{"Field", "Value"}, rows)
		},
	}

	cmd.Flags().StringVar(&projectID, "project", "", "Project ID")
	cmd.Flags().StringVar(&taskID, "id", "", "Task ID")
	cmd.Flags().StringVar(&format, "format", "", "Output format: table, json")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}
