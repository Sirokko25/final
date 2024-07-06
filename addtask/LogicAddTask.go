package addtask

import (
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"final/constant"
	"final/nextdate"
	"final/task"
)

func Checktitle(task task.Task) error {
	if task.Title == "" {
		return errors.New("Пустой заголовок")
	}
	return nil
}
func Checkdate(task task.Task) (task.Task, error) {

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(constant.ParseDate)
		return task, nil
	} else {
		date, err := time.Parse(constant.ParseDate, task.Date)
		if err != nil {
			return task, errors.New("Неправильный формат даты")
		}
		if date.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(constant.ParseDate)
				return task, nil
			} else {
				nowtime := now.Format(constant.ParseDate)
				if nowtime != task.Date {
					nextDate, err := nextdate.CalcNextDate(nowtime, task.Date, task.Repeat)
					if err != nil {
						return task, errors.New("Ошибка вычисления даты")
					}
					task.Date = nextDate
					return task, nil
				} else {
					task.Date = nowtime
					return task, nil
				}
			}
		}
	}
	return task, nil
}

func Checkrulerepeat(task task.Task) error {
	if task.Repeat != "" {
		now := time.Now()
		nowtime := now.Format(constant.ParseDate)
		nextDate, err := nextdate.CalcNextDate(nowtime, task.Date, task.Repeat)
		if err != nil {
			return errors.New("Ошибка вычисления даты")
		}
		task.Date = nextDate
	}
	return nil
}
