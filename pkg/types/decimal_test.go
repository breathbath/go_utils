package types

import (
	"testing"

	coreDecimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDecimalFromStringCreation(t *testing.T) {
	assertEqualDecimals(t, "3.14", NewDecimalFromString("3.14"))

	assertEqualDecimals(t, "0", NewDecimalFromString("invalidDecimal"))
}

func TestDecimalFromFloatCreation(t *testing.T) {
	assertEqualDecimals(t, "3.14", NewDecimalFromFloat(3.14))
}

func TestDecimalFromIntCreation(t *testing.T) {
	assertEqualDecimals(t, "314", NewDecimalFromInt(314))
}

func TestDecimalToString(t *testing.T) {
	actualDecimal := NewDecimalFromString("3.12345678901234567890")
	assert.Equal(t, "3.1234567890123", actualDecimal.String())
}

func TestDecimalToFloat(t *testing.T) {
	actualDecimal := NewDecimalFromString("3.123")
	assert.Equal(t, 3.123, actualDecimal.ToFloat())
}

func TestDecimalMathOperations(t *testing.T) {
	decimalOne := NewDecimalFromString("3.123")
	decimalTwo := NewDecimalFromString("0.1")
	decimalZero := NewDecimalFromString("0")

	assertEqualDecimals(t, "3.223", decimalOne.Add(decimalTwo))

	assertEqualDecimals(t, "3.023", decimalOne.Sub(decimalTwo))

	assertEqualDecimals(t, "31.23", decimalOne.Div(decimalTwo))

	assertEqualDecimals(t, "0.3123", decimalOne.Mul(decimalTwo))

	assertEqualDecimals(t, "3.2", decimalOne.CeilToValue(decimalTwo))
	assertEqualDecimals(t, "0", decimalZero.CeilToValue(NewDecimalFromInt(0)))

	assertEqualDecimals(t, "3.1", decimalOne.FloorToValue(decimalTwo))
	assertEqualDecimals(t, "0", decimalZero.FloorToValue(NewDecimalFromInt(0)))

	assertEqualDecimals(t, "3", decimalOne.Floor())

	assertEqualDecimals(t, "3.4353", decimalOne.IncrementByPercent(decimalTwo))
	assertEqualDecimals(t, "2.8107", decimalOne.DecrementByPercent(decimalTwo))

	assertEqualDecimals(t, "3.123", decimalOne)
	assertEqualDecimals(t, "0.1", decimalTwo)
}

func TestJsonConversions(t *testing.T) {
	decimalVal := NewDecimalFromString("3.123")
	jsonDecimal, err := decimalVal.MarshalJSON()

	assert.NoError(t, err)
	assert.Equal(t, "3.123", string(jsonDecimal))

	err = decimalVal.UnmarshalJSON([]byte("1.2"))
	assert.NoError(t, err)
	assertEqualDecimals(t, "1.2", decimalVal)
}

func TestDecimalScan(t *testing.T) {
	decimalVal := Decimal{}

	err := decimalVal.Scan(nil)
	assert.NoError(t, err)
	assertEqualDecimals(t, "0", decimalVal)

	err = decimalVal.Scan(0.2)
	assert.NoError(t, err)
	assertEqualDecimals(t, "0.2", decimalVal)

	err = decimalVal.Scan(2.0)
	assert.NoError(t, err)
	assertEqualDecimals(t, "2", decimalVal)

	err = decimalVal.Scan(float64(2))
	assert.NoError(t, err)
	assertEqualDecimals(t, "2", decimalVal)

	err = decimalVal.Scan(float32(4))
	assert.NoError(t, err)
	assertEqualDecimals(t, "4", decimalVal)
}

func TestValue(t *testing.T) {
	decimalVal := NewDecimalFromString("3.123")
	val, err := decimalVal.Value()
	assert.NoError(t, err)

	assert.Equal(t, "3.123", val)
}

func TestComparing(t *testing.T) {
	decimalOne := NewDecimalFromString("3.123")
	decimalTwo := NewDecimalFromString("0.1")
	decimalThree := NewDecimalFromString("3.123")
	decimalFour := NewDecimalFromString("-3.123")
	decimalFive := NewDecimalFromString("5")

	assert.True(t, decimalTwo.LessOrEqual(decimalOne))
	assert.True(t, decimalTwo.LessOrEqual(decimalThree))
	assert.True(t, decimalOne.Equal(decimalThree))
	assert.True(t, decimalOne.GreaterOrEqual(decimalTwo))
	assert.True(t, decimalThree.GreaterOrEqual(decimalOne))
	assert.True(t, decimalTwo.Less(decimalOne))
	assert.True(t, decimalOne.Greater(decimalTwo))
	assert.True(t, decimalFour.LessOrEqualZero())
	assert.True(t, ZERO.LessOrEqualZero())
	assert.True(t, ZERO.GreaterOrEqualZero())
	assert.True(t, decimalOne.GreaterZero())
	assert.True(t, decimalFour.LessZero())
	assert.False(t, ZERO.GreaterZero())
	assert.False(t, ZERO.LessZero())
	assert.True(t, decimalFive.EqualInt(5))
	assert.True(t, decimalOne.GreaterInt(3))
	assert.True(t, decimalFive.LowerOrEqualInt(5))
	assert.True(t, decimalOne.LowerOrEqualInt(4))
	assert.True(t, decimalFive.GreaterOrEqualInt(5))
	assert.True(t, decimalOne.GreaterOrEqualInt(3))
}

type RoundingTestSet struct {
	inputNumber    string
	places         int64
	expectedResult string
}

func TestToPercentConversion(t *testing.T) {
	testingSets := []RoundingTestSet{
		{inputNumber: "0.123", places: 2, expectedResult: "12.3"},
		{inputNumber: "1.123", places: 2, expectedResult: "112.3"},
		{inputNumber: "0.123456", places: -1, expectedResult: "12.3456"},
		{inputNumber: "0.123456", places: 2, expectedResult: "12.35"},
		{inputNumber: "0.123333", places: 2, expectedResult: "12.33"},
		{inputNumber: "1", places: 2, expectedResult: "100"},
		{inputNumber: "10", places: 2, expectedResult: "1000"},
		{inputNumber: "-0.991", places: 2, expectedResult: "-99.1"},
		{inputNumber: "-0.9999", places: 0, expectedResult: "-100"},
	}

	for _, testingSet := range testingSets {
		decimalVal := NewDecimalFromString(testingSet.inputNumber)
		convertedDecimalVal := decimalVal.ToPercent(testingSet.places)
		assertEqualDecimals(t, testingSet.expectedResult, convertedDecimalVal)
	}
}

func TestRounding(t *testing.T) {
	testingSets := []RoundingTestSet{
		{inputNumber: "0.123", places: 2, expectedResult: "0.12"},
		{inputNumber: "1.125", places: 2, expectedResult: "1.13"},
		{inputNumber: "0.123456", places: -1, expectedResult: "0.123456"},
		{inputNumber: "1.1", places: 0, expectedResult: "1"},
		{inputNumber: "1.9", places: 0, expectedResult: "2"},
		{inputNumber: "1", places: 2, expectedResult: "1"},
		{inputNumber: "0.12", places: 2, expectedResult: "0.12"},
		{inputNumber: "-0.991", places: 2, expectedResult: "-0.99"},
		{inputNumber: "-0.9999", places: 3, expectedResult: "-1"},
	}

	for _, testingSet := range testingSets {
		decimalVal := NewDecimalFromString(testingSet.inputNumber)
		convertedDecimalVal := decimalVal.Round(testingSet.places)
		assertEqualDecimals(t, testingSet.expectedResult, convertedDecimalVal)
	}
}

func assertEqualDecimals(t *testing.T, expectedDecimalStr string, actualDecimal Decimal) {
	coreDec, err := coreDecimal.NewFromString(expectedDecimalStr)
	assert.NoError(t, err)
	expectedDecimal := Decimal{dec: coreDec}

	assert.Equal(t, expectedDecimal.String(), actualDecimal.String())
}
