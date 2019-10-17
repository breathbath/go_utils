package testing

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCaptureOutput(t *testing.T) {
	outputFunc := func() {
		log.Println("Some text")
	}

	actualOutput := CaptureOutput(outputFunc)
	assert.Contains(t, actualOutput, "Some text")
}
