package math

import (
	"fmt"
	"github.com/breathbath/go_utils/utils/conv"
	"math"
	"math/rand"
	"strings"
	"time"
)

func Round(f float64, places float64) float64 {
	shift := math.Pow(10, places)
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

func CountDigits(input int64) (int64) {
	strNumb := fmt.Sprintf("%d", input)
	strLen := int64(len(strNumb))
	if input < 0 {
		return strLen - 1
	}

	return strLen
}
