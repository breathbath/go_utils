package conv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockStruct struct{}

func (ms MockStruct) String() string {
	return "mockStruct"
}

func TestStringerToJson(t *testing.T) {
	ms := MockStruct{}
	output := StringerToJson(ms)
	assert.Equal(t, `"mockStruct"`, string(output))
}
