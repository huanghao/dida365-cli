package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(app *App) *cobra.Command {
	var projectID string
	var taskID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(projectID) == "" || strings.TrimSpace(taskID) == "" {
				return fmt.Errorf("--project and --id are required")
			}
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			if app.DryRun {
				fmt.Fprintf(app.Out, "Would call DELETE %s/project/%s/task/%s\n", cfg.APIBaseURL, projectID, taskID)
				return nil
			}
			client := newAPIClient(cfg)
			if err := client.DeleteTask(projectID, taskID); err != nil {
				return err
			}
			fmt.Fprintf(app.Out, "Deleted task: %s\n", taskID)
			return nil
		},
	}

	cmd.Flags().StringVar(&projectID, "project", "", "Project ID")
	cmd.Flags().StringVar(&taskID, "id", "", "Task ID")
	return cmd
}
