package utils

func IntArrayContain(array []int, value int) bool {
	for _, item := range array {
		if value == item {
			return true
		}
	}
	return false
}

func Int64ArrayContain(array []int64, value int64) bool {
	for _, item := range array {
		if value == item {
			return true
		}
	}
	return false
}
