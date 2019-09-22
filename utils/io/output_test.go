package io

import (
	"errors"
	testing2 "github.com/breathbath/go_utils/utils/testing"
	"github.com/stretchr/testify/assert"
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

func TestOutputError(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		err := errors.New("Some error")
		OutputError(err, "Some topic", "Some msg %s", "text")
	})
	testing2.AssertLogText(t,"[ERROR] Some error, Some msg text [Some topic]", output)
}

func TestOutputWarning(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputWarning("Topic", "Number %d", 1)
	})
	testing2.AssertLogText(t,"[WARNING] Number 1 [Topic]", output)
}

func TestOutputInfo(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputInfo("Info_top", "Many params %d, %s", 1, "lala")
	})
	testing2.AssertLogText(t,"[INFO] Many params 1, lala [Info_top]", output)
}

func TestOutputWithFormatChars(t *testing.T) {
	a := "Some msg 10%---20%s"
	output := testing2.CaptureOutput(func() {
		OutputMessageType("INFO", "Top",  a)
	})
	testing2.AssertLogText(t,"[INFO] Some msg 10%---20%s [Top]", output)
}

func TestOutputMsgType(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputMessageType("Type", "Topic", "Msg")
	})
	testing2.AssertLogText(t,"[Type] Msg [Topic]", output)
}

func TestOutputMsgWithoutTopic(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputMessageType("Type", "", "Msg")
	})
	testing2.AssertLogText(t,"[Type] Msg", output)
}

func TestCutMessageWithNoLimit(t *testing.T) {
	SetMaxMessageLength(0)
	msg := CutMessageIfNeeded("some msg")
	assert.Equal(t, "some msg", msg)
}
