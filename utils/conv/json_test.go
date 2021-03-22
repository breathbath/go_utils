package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStruct struct{}

func (ms MockStruct) String() string {
	return "mockStruct"
}

func TestStringerToJson(t *testing.T) {
	ms := MockStruct{}
	output := StringerToJSON(ms)
	assert.Equal(t, `"mockStruct"`, string(output))
}
