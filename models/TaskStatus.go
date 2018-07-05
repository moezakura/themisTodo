package models

import (
	"strings"
	"errors"
)

type TaskStatus int

const (
	TASK_STATUS_TODO         TaskStatus = iota
	TASK_STATUS_DOING
	TASK_STATUS_PULL_REQUEST
	TASK_STATUS_DONE
	TASK_STATUS_HIDE
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
	case TASK_STATUS_HIDE:
		return "HIDE"
	case TASK_STATUS_OTHER:
		return "OTHER"
	default:
		return "Unknown"
	}
}

func TaskStatusFromString(statusText string) (status TaskStatus, err error) {
	statusText = strings.ToUpper(statusText)

	switch statusText {
	case "TODO":
		return TASK_STATUS_TODO, nil
	case "DOING":
		return TASK_STATUS_DOING, nil
	case "PULL_REQUEST":
	case "PULLREQUEST":
		return TASK_STATUS_PULL_REQUEST, nil
	case "DONE":
		return TASK_STATUS_DONE, nil
	case "HIDE":
		return TASK_STATUS_HIDE, nil
	}

	return TASK_STATUS_OTHER, errors.New("invalid status.")
}
