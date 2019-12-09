package utils

import (
	"themis.mox.si/themis/models"
	"log"
)

func TasksConvert(tasks []models.Task) []models.Task {
	for key, value := range tasks {
		tasks[key] = *TaskConvert(&value)
	}

	return tasks
}

func TaskConvert(task *models.Task) *models.Task {
	var e bool
	e, task.DeadlineMD = GetDateMD(task.Deadline)
	if e {
		log.Printf("Utils.TaskConvert")
	}

	task.LimitDate = DiffDay(task.Deadline)

	return task
}

func TaskHistoryItemConvert(task *models.TaskHistoryItem) *models.TaskHistoryItem {
	var e bool
	e, task.DeadlineMD = GetDateMD(task.Deadline)
	if e {
		log.Printf("Utils.TaskConvert")
	}

	task.LimitDate = DiffDay(task.Deadline)

	return task
}
