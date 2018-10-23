package types

import (
	coreDecimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
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

	assertEqualDecimals(t, "3.223", decimalOne.Add(decimalTwo))

	assertEqualDecimals(t, "3.023", decimalOne.Sub(decimalTwo))

	assertEqualDecimals(t, "31.23", decimalOne.Div(decimalTwo))

	assertEqualDecimals(t, "0.3123", decimalOne.Mul(decimalTwo))

	assertEqualDecimals(t, "3.2", decimalOne.CeilToValue(decimalTwo))

	assertEqualDecimals(t, "3.1", decimalOne.FloorToValue(decimalTwo))

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

func assertEqualDecimals(t *testing.T, expectedDecimalStr string, actualDecimal Decimal) {
	coreDec, err := coreDecimal.NewFromString(expectedDecimalStr)
	assert.NoError(t, err)
	expectedDecimal := Decimal{dec: coreDec}

	assert.Equal(t, expectedDecimal.String(), actualDecimal.String())
}
