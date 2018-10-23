package io

import (
	testing2 "github.com/breathbath/go_utils/utils/testing"
	"testing"
)

const(
	MULTILINE_TEXT=`One
per
each
line`
)

func TestOutputSingleLine(t *testing.T) {
	SetMaxMessageLength(100)
	output := testing2.CaptureOutput(func() {
		OutputSingleLine("Some msg %s", "text")
	})

	testing2.AssertLogText(t, "Some msg text", output)

	output = testing2.CaptureOutput(func() {
		OutputSingleLine(MULTILINE_TEXT)
	})

	testing2.AssertLogText(t,"One per each line", output)

	SetMaxMessageLength(10)
	output = testing2.CaptureOutput(func() {
		OutputSingleLine("Too long string to show")
	})

	testing2.AssertLogText(t,"Too long s", output)
}

func TestOutputSingleLineWithTopic(t *testing.T) {
	SetMaxMessageLength(100)
	output := testing2.CaptureOutput(func() {
		OutputSingleLineWithTopic("Some topic", "Some msg %s", "text")
	})
	testing2.AssertLogText(t,"[Some topic] Some msg text", output)
}