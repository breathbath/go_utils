package time

import "time"

func GetDateAtDayStart(inputDate time.Time) time.Time {
	return time.Date(inputDate.Year(), inputDate.Month(), inputDate.Day(), 0, 0, 0, 0, inputDate.Location())
}

func GetDateAtDayEnd(inputDate time.Time) time.Time {
	return time.Date(inputDate.Year(), inputDate.Month(), inputDate.Day(), 23, 59, 59, 999999999, inputDate.Location())
}
