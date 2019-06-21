package models

type TaskTimerUpdateRequestJson struct {
	StartDateHMS string `json:"start_date_hms"`
	EndDateHMS   string `json:"end_date_hms"`
	Note        string `json:"note"`
}
