package models

import "time"

type TodoTimer struct {
	Id            int       `json:"id"`
	CreateDate    int64     `json:"create_date"`
	Assign        int       `json:"assign"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	StartDateUnix int64     `json:"start_date_unix"`
	EndDateUnix   int64     `json:"end_date_unix"`
	Note          []byte
	NoteString    string `json:"note"`
}
