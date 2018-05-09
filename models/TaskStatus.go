package models

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