package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleStringNumber = StringNumber{"123"}

func TestStringNumberMarshalJSON(t *testing.T) {
	byteStringNumber, err := exampleStringNumber.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "123", string(byteStringNumber))

	nilStringNumber := StringNumber{}
	byteStringNumber, err = nilStringNumber.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "null", string(byteStringNumber))
}

func TestStringNumberUnMarshalJSON(t *testing.T) {
	someStringNumber := StringNumber{}

	err := someStringNumber.UnmarshalJSON([]byte("123"))
	assert.NoError(t, err)
	assert.Equal(t, StringNumber{"123"}, someStringNumber)

	err = someStringNumber.UnmarshalJSON([]byte(nil))
	assert.NoError(t, err)
	assert.Equal(t, StringNumber{}, someStringNumber)

	err = someStringNumber.UnmarshalJSON([]byte("not_number"))
	assert.EqualError(t, err, "cannot unmarshal 'not_number' into a number type")
}

func TestStringNumberStringer(t *testing.T) {
	someStringNumber := StringNumber{}
	assert.Equal(t, "", someStringNumber.String())

	someStringNumber = StringNumber{"222"}
	assert.Equal(t, "222", someStringNumber.String())
}
