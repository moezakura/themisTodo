package models

import "time"

type TodoTimer struct {
	Id               int         `json:"id"`
	CreateDate       int64       `json:"-"`
	CreateDateString string      `json:"create_date"`
	Assign           int         `json:"assign"`
	StartDate        time.Time   `json:"-"`
	EndDate          time.Time   `json:"-"`
	StartDateUnix    int64       `json:"start_date_unix"`
	EndDateUnix      int64       `json:"end_date_unix"`
	Note             []byte      `json:"-"`
	NoteString       string      `json:"note"`
	Task             *TaskOfJson `json:"task"`
}
