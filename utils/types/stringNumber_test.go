package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var exampleStringNumber StringNumber

func init() {
	exampleStringNumber = StringNumber{"123"}
}

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
	assert.EqualError(t, err, "Cannot unmarshal 'not_number' into a number type")
}
