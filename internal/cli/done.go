package cli

import (
	"fmt"
	"strings"

	"github.com/huanghao/dida365-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewDoneCommand(app *App) *cobra.Command {
	var projectID string
	var taskID string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "done",
		Short: "Complete a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(projectID) == "" || strings.TrimSpace(taskID) == "" {
				return fmt.Errorf("--project and --id are required")
			}
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			if app.DryRun {
				if asJSON {
					return output.PrintJSON(app.Out, map[string]any{
						"action":  "complete_task",
						"dry_run": true,
						"method":  "POST",
						"url":     fmt.Sprintf("%s/project/%s/task/%s/complete", cfg.APIBaseURL, projectID, taskID),
						"project": projectID,
						"task_id": taskID,
					})
				}
				fmt.Fprintf(app.Out, "Would call POST %s/project/%s/task/%s/complete\n", cfg.APIBaseURL, projectID, taskID)
				return nil
			}
			client := newAPIClient(cfg)
			if err := client.CompleteTask(projectID, taskID); err != nil {
				return err
			}
			if asJSON {
				return output.PrintJSON(app.Out, map[string]any{
					"ok":      true,
					"action":  "complete_task",
					"project": projectID,
					"task_id": taskID,
				})
			}
			fmt.Fprintf(app.Out, "Completed task: %s\n", taskID)
			return nil
		},
	}

	cmd.Flags().StringVar(&projectID, "project", "", "Project ID")
	cmd.Flags().StringVar(&taskID, "id", "", "Task ID")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}
