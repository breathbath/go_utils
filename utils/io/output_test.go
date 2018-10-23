package io

import (
	"fmt"
	testing2 "github.com/breathbath/go_utils/utils/testing"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

const(
	MULTILINE_TEXT=`One
per
each
line`
)

func TestOutputSingleLine(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputSingleLine("Some msg %s", "text")
	})

	AssertLogText(t, "Some msg text", output)

	output = testing2.CaptureOutput(func() {
		OutputSingleLine(MULTILINE_TEXT)
	})

	AssertLogText(t,"One per each line", output)

	SetMaxMessageLength(10)
	output = testing2.CaptureOutput(func() {
		OutputSingleLine("Too long string to show")
	})

	AssertLogText(t,"Too long s", output)
}

func TestOutputSingleLineWithTopic(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputSingleLineWithTopic("Some topic", "Some msg %s", "text")
	})
	AssertLogText(t,"[Some topic] Some msg text", output)
}

func AssertLogText(t *testing.T, expectedText, actualText string) {
	regx := fmt.Sprintf(`\d{4}\/\d{2}\/\d{2}\s\d{2}:\d{2}:\d{2}\s%s`, regexp.QuoteMeta(expectedText))
	assert.Regexp(t, regx, actualText)
}