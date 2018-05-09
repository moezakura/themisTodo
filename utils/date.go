package utils

import "time"

func ParseDate(dateText string) (e bool, t time.Time) {
	userTime, err := time.Parse("2006-01-02", dateText)
	if err != nil {
		return true, time.Time{}
	}

	return false, userTime
}

func GetDateMD(dateText string) (e bool, date string) {
	err, userTime := ParseDate(dateText)
	if err {
		return true, ""
	}

	return false, userTime.Format("01/02")
}

func DiffDay(dateText string) int {
	err, userTime := ParseDate(dateText)
	if err {
		return 0
	}

	diff := userTime.Sub(time.Now())

	return int(diff.Hours() / 24)
}
