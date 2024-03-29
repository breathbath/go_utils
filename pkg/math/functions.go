package math

import (
	"fmt"
	"math"
	"strings"

	"github.com/breathbath/go_utils/v3/pkg/conv"
)

const DefaultRoundingBase = 10
const DefaultRoundingPlaces = .5

func Round(f, places float64) float64 {
	shift := math.Pow(DefaultRoundingBase, places)
	return math.Floor(f*shift+DefaultRoundingPlaces) / shift
}

func CountDecimalPlaces(input float64) int {
	strNumb := conv.ConvertFloatToLongStringNumber(input)
	splitNumber := strings.Split(strNumb, ".")
	if len(splitNumber) == 1 {
		return 0
	}
	return len([]rune(splitNumber[1]))
}

func CountDigits(input int64) int64 {
	strNumb := fmt.Sprintf("%d", input)
	strLen := int64(len(strNumb))
	if input < 0 {
		return strLen - 1
	}

	return strLen
}
