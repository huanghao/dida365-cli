package cli

import (
	"fmt"
	"strings"

	"github.com/huanghao/dida365-cli/internal/dida"
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
	cmd.AddCommand(newProjectsCreateCommand(app))
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
			cacheKey := readCacheKey("GET", "/project")
			var projects []dida.Project
			if !readCacheGet(app, cacheKey, &projects) {
				projects, err = client.GetProjects()
				if err != nil {
					return err
				}
				readCachePut(app, cacheKey, projects)
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

func newProjectsCreateCommand(app *App) *cobra.Command {
	var name string
	var color string
	var sortOrder int64
	var hasSortOrder bool
	var viewMode string
	var kind string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf("--name is required")
			}
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			input := dida.CreateProjectInput{
				Name:     name,
				Color:    color,
				ViewMode: viewMode,
				Kind:     kind,
			}
			if hasSortOrder {
				input.SortOrder = &sortOrder
			}
			if app.DryRun {
				fmt.Fprintf(app.Out, "Would call POST %s/project\n", cfg.APIBaseURL)
				return output.PrintJSON(app.Out, input)
			}
			if err := checkWriteDebounce(app, "create_project", input); err != nil {
				return err
			}
			client := newAPIClient(cfg)
			project, err := client.CreateProject(input)
			if err != nil {
				return err
			}
			markWriteDebounce(app, "create_project", input)
			clearReadCache(app)
			if asJSON {
				return output.PrintJSON(app.Out, project)
			}
			rows := [][]string{
				{"ID", project.ID},
				{"Name", project.Name},
				{"Kind", project.Kind},
				{"ViewMode", project.ViewMode},
				{"Color", project.Color},
				{"SortOrder", fmt.Sprintf("%d", project.SortOrder)},
			}
			return output.PrintSimpleTable(app.Out, []string{"Field", "Value"}, rows)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Project name (required)")
	cmd.Flags().StringVar(&color, "color", "", "Project color, e.g. #F18181")
	cmd.Flags().Int64Var(&sortOrder, "sort-order", 0, "Project sort order")
	cmd.Flags().StringVar(&viewMode, "view-mode", "", "Project view mode: list|kanban|timeline")
	cmd.Flags().StringVar(&kind, "kind", "", "Project kind: TASK|NOTE")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		hasSortOrder = cmd.Flags().Changed("sort-order")
	}

	return cmd
}
