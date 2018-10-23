package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	someDecimal := NewDecimalFromString("2.12")

	nullDecimal := NullDecimal{
		DecimalValue: someDecimal,
		Valid:        false,
	}

	jsonDecimal, err := nullDecimal.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "null", string(jsonDecimal))

	nullDecimal.Valid = true
	jsonDecimal, err = nullDecimal.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `2.12`, string(jsonDecimal))
}

func TestNullDecimalScan(t *testing.T) {
	someNullDecimal := NullDecimal{}
	invalidValues := []interface{}{
		nil,
		[]byte(""),
		"",
	}

	for _, invalidValue := range invalidValues {
		err := someNullDecimal.Scan(invalidValue)
		assert.NoError(t, err)
		assert.Equal(t, ZERO, someNullDecimal.DecimalValue)
		assert.False(t, someNullDecimal.Valid)
	}

	validValues := []interface{}{
		3.14,
		[]byte("3.14"),
		"3.14",
	}

	expectedDecimal := NewDecimalFromString("3.14")
	for _, validValue := range validValues {
		err := someNullDecimal.Scan(validValue)
		assert.NoError(t, err)
		assert.True(t, someNullDecimal.Valid)
		assert.Equal(t, expectedDecimal, someNullDecimal.DecimalValue)
	}
}
