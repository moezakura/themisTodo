package models

type TodoTimer struct {
	Id         int    `json:"id"`
	CreateDate int64  `json:"create_date"`
	Assign     int    `json:"assign"`
	StartDate  int64  `json:"start_date"`
	EndDate    int64  `json:"end_date"`
	Note       string `json:"note"`
}
