package cli

import (
	"fmt"
	"unicode/utf8"

	"github.com/huanghao/dida365-cli/internal/dida"
)

const maxCreateTextRunes = 499

func validateCreateTaskInput(title, content, desc string) error {
	if err := validateMaxRunes("title", title, maxCreateTextRunes); err != nil {
		return err
	}
	if err := validateMaxRunes("content", content, maxCreateTextRunes); err != nil {
		return err
	}
	if err := validateMaxRunes("desc", desc, maxCreateTextRunes); err != nil {
		return err
	}
	return nil
}

func validateMaxRunes(field, value string, max int) error {
	if utf8.RuneCountInString(value) > max {
		return fmt.Errorf("--%s length must be < 500 characters", field)
	}
	return nil
}

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
