package testing

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureOutput(t *testing.T) {
	outputFunc := func() {
		log.Println("Some text")
	}

	actualOutput := CaptureOutput(outputFunc)
	assert.Contains(t, actualOutput, "Some text")
}
