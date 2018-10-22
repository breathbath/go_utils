package math

import (
	"github.com/breathbath/go_utils/utils/conv"
	"math"
	"math/rand"
	"strings"
	"time"
)

func Round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

func CountDecimalPlaces(input float64) int {
	strNumb := conv.ConvertFloatToLongStringNumber(input)
	splitNumber := strings.Split(strNumb, ".")
	if len(splitNumber) == 1 {
		return 0
	}
	return len([]rune(splitNumber[1]))
}

func RandInt(limit int64) int64 {
	limitInt := int(limit)
	rand.Seed(time.Now().UnixNano())

	return int64(rand.Intn(limitInt))
}
