package changetask

import (
	"time"

	_ "github.com/mattn/go-sqlite3"

	"final/constant"
	"final/nextdate"
	"final/task"

)

func CheckId(task task.Task) string {
	if task.ID == "" {
		return "Не указан индентификатор задачи"
	} else {
		return ""
	}
}

func CheckTitle(task task.Task) string {
	if task.Title == "" {
		return "Не указан заголовок задачи"
	} else {
		return ""
	}
}

func CheckDate(task task.Task) (task.Task, string) {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(constant.ParseDate)
		return task, ""
	} else {
		date, err := time.Parse(constant.ParseDate, task.Date)
		if err != nil {
			return task, "Неправильный формат даты"
		}
		if date.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(constant.ParseDate)
				return task, ""
			} else {
				nowtime := now.Format(constant.ParseDate)
				nextDate, err := nextdate.CalcNextDate(nowtime, task.Date, task.Repeat)
				if err != nil {
					return task, "Ошибка вычисления даты"
				}
				task.Date = nextDate
				return task, ""
			}

		}
		return task, ""
	}
}

func CheckRepeate(task task.Task) string {
	if task.Repeat != "" {
		_, err := nextdate.ParseRepeatRules(task.Repeat)
		if err != nil {
			return "Правило повторения указано в неправильном формате"
		}
	}
	return ""
}
