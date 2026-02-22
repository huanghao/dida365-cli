package cli

import (
	"fmt"
	"strings"

	"github.com/huanghao/dida365-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewAddCommand(app *App) *cobra.Command {
	var projectID string
	var title string
	var content string
	var desc string
	var startDate string
	var dueDate string
	var repeatFlag string
	var timeZone string
	var priority int
	var allDay bool
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Create a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(projectID) == "" || strings.TrimSpace(title) == "" {
				return fmt.Errorf("--project and --title are required")
			}

			input := buildTaskFromFlags(projectID, title, content, desc, startDate, dueDate, repeatFlag, timeZone, allDay, priority)
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			if app.DryRun {
				fmt.Fprintf(app.Out, "Would call POST %s/task\n", cfg.APIBaseURL)
				return output.PrintJSON(app.Out, input)
			}

			client := newAPIClient(cfg)
			task, err := client.CreateTask(input)
			if err != nil {
				return err
			}
			if asJSON {
				return output.PrintJSON(app.Out, task)
			}
			fmt.Fprintf(app.Out, "Created task: %s (%s)\n", task.Title, task.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&projectID, "project", "", "Project ID")
	cmd.Flags().StringVar(&title, "title", "", "Task title")
	cmd.Flags().StringVar(&content, "content", "", "Task content")
	cmd.Flags().StringVar(&desc, "desc", "", "Task description")
	cmd.Flags().StringVar(&startDate, "start", "", "Start datetime (yyyy-MM-dd'T'HH:mm:ssZ)")
	cmd.Flags().StringVar(&dueDate, "due", "", "Due datetime (yyyy-MM-dd'T'HH:mm:ssZ)")
	cmd.Flags().StringVar(&repeatFlag, "repeat", "", "Repeat rule, e.g. RRULE:FREQ=DAILY;INTERVAL=1")
	cmd.Flags().StringVar(&timeZone, "timezone", "", "Time zone, e.g. America/Los_Angeles")
	cmd.Flags().IntVar(&priority, "priority", 0, "Priority (default 0)")
	cmd.Flags().BoolVar(&allDay, "all-day", false, "Set all-day task")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")

	return cmd
}
