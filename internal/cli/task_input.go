package cli

import "github.com/huanghao/dida365-cli/internal/dida"

func buildTaskFromFlags(projectID, title, content, desc, startDate, dueDate, repeatFlag, timeZone string, allDay bool, priority int) dida.Task {
	return dida.Task{
		ProjectID:  projectID,
		Title:      title,
		Content:    content,
		Desc:       desc,
		StartDate:  startDate,
		DueDate:    dueDate,
		RepeatFlag: repeatFlag,
		TimeZone:   timeZone,
		IsAllDay:   allDay,
		Priority:   priority,
	}
}
