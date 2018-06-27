package models

import "github.com/pkg/errors"

type TaskStatus int

const (
	TASK_STATUS_TODO TaskStatus = iota
	TASK_STATUS_DOING
	TASK_STATUS_PULL_REQUEST
	TASK_STATUS_DONE
	TASK_STATUS_OTHER
)

func (rm TaskStatus) String() string {
	switch rm {
	case TASK_STATUS_TODO:
		return "TODO"
	case TASK_STATUS_DOING:
		return "DOING"
	case TASK_STATUS_PULL_REQUEST:
		return "PULL_REQUEST"
	case TASK_STATUS_DONE:
		return "DONE"
	case TASK_STATUS_OTHER:
		return "OTHER"
	default:
		return "Unknown"
	}
}

func StringToTaskStatus(rm string) (TaskStatus, error) {
	switch rm {
	case "TODO", "todo":
		return TASK_STATUS_TODO, nil
	case "DOING", "doing":
		return TASK_STATUS_DOING, nil
	case "PULL_REQUEST", "pull_request":
		return TASK_STATUS_PULL_REQUEST, nil
	case "DONE", "done":
		return TASK_STATUS_DONE, nil
	case "OTHER", "other":
		return TASK_STATUS_OTHER, nil
	default:
		return -1, errors.New("Unknown task status")
	}
}