package utils

import "time"

func ParseDateOrNow(dateStr string) time.Time {
	layout := "02-01-2006" // dd-MM-yyyy
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Now()
	}
	return t
}
