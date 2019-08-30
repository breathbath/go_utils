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

func TestJsonUnMarshal(t *testing.T) {
	jsonDecimal := `2.34`
	nd := NullDecimal{}

	err := nd.UnmarshalJSON([]byte(jsonDecimal))
	assert.NoError(t, err)

	expectedDecimal := NewDecimalFromString("2.34")

	assert.True(t, nd.Valid)
	assert.Equal(t, expectedDecimal, nd.DecimalValue)

	err = nd.UnmarshalJSON([]byte("invalid decimal"))
	assert.EqualError(t, err, "Error decoding string 'invalid decimal': can't convert invalid decimal to decimal: exponent is not numeric")

	err = nd.UnmarshalJSON([]byte(`"22"`))
	expectedDecimal = NewDecimalFromString("22")
	assert.True(t, nd.Valid)
	assert.Equal(t, expectedDecimal, nd.DecimalValue)
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
		driverVal, err := someNullDecimal.Value()
		assert.NoError(t, err)
		assert.NotNil(t, driverVal)
	}
}

func TestInvalidNullDecimal(t *testing.T) {
	invalidNullDecimal := NullDecimal{Valid: false}
	val, err := invalidNullDecimal.Value()
	assert.Nil(t, val)
	assert.NoError(t, err)
}
