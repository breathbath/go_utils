package testing

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertLogText(t *testing.T) {
	outputFunc := func() {
		log.Println("Some text")
	}

	actualOutput := CaptureOutput(outputFunc)

	localT := &testing.T{}
	AssertLogText(localT, "Some text", actualOutput)

	isFailed := localT.Failed()
	assert.False(t, isFailed)
}
