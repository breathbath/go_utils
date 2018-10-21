package conv

import (
	"regexp"
	"strconv"
)

func ExtractIntFromString(input string, defaultVal int64) int64 {
	re := regexp.MustCompile(`\D*`)
	strInt := re.ReplaceAllString(input, "")

	intVal, err := strconv.ParseInt(strInt, 10, 64)
	if err != nil {
		return defaultVal
	}

	return intVal
}

func ConvertFloatToLongStringNumber(input float64) (string) {
	return strconv.FormatFloat(input, 'f', -1, 64)
}
