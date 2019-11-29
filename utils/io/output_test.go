package io

import (
	"errors"
	"fmt"
	testing2 "github.com/breathbath/go_utils/utils/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	MULTILINE_TEXT = `One
per
each
line`
)

func TestOutputSingleLine(t *testing.T) {
	SetMaxMessageLength(100)
	output := testing2.CaptureOutput(func() {
		OutputInfo("", "Some msg %s", "text")
	})

	testing2.AssertLogText(t, fmt.Sprintf("[%s] Some msg text", SeverityInfo), output)

	output = testing2.CaptureOutput(func() {
		OutputInfo("", MULTILINE_TEXT)
	})

	testing2.AssertLogText(t, fmt.Sprintf("[%s] One per each line", SeverityInfo), output)

	SetMaxMessageLength(10)
	output = testing2.CaptureOutput(func() {
		OutputInfo("", "Too long string to show")
	})

	testing2.AssertLogText(t, fmt.Sprintf("[%s] Too long s", SeverityInfo), output)
}

func TestOutputSingleLineWithTopic(t *testing.T) {
	SetMaxMessageLength(100)
	output := testing2.CaptureOutput(func() {
		OutputInfo("Some topic", "Some msg %s", "text")
	})
	testing2.AssertLogText(t, fmt.Sprintf("[%s] Some msg text [Some topic]", SeverityInfo), output)
}

func TestOutputError(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		err := errors.New("Some error")
		OutputError(err, "Some topic", "Some msg %s", "text")
	})
	testing2.AssertLogText(t, fmt.Sprintf("[%s] Some error, Some msg text [Some topic]", SeverityError), output)
}

func TestOutputWarning(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputWarning("Topic", "Number %d", 1)
	})
	testing2.AssertLogText(t, fmt.Sprintf("[%s] Number 1 [Topic]", SeverityWarning), output)
}

func TestOutputInfo(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputInfo("Info_top", "Many params %d, %s", 1, "lala")
	})
	testing2.AssertLogText(t, fmt.Sprintf("[%s] Many params 1, lala [Info_top]", SeverityInfo), output)
}

func TestOutputWithFormatChars(t *testing.T) {
	a := "Some msg 10%---20%s"
	output := testing2.CaptureOutput(func() {
		OutputMessageType(SeverityInfo, "Top", a)
	})
	testing2.AssertLogText(t, "[" + SeverityInfo + "] Some msg 10%---20%s [Top]", output)

	b := "Some msg 10%%---20%s"
	output2 := testing2.CaptureOutput(func() {
		OutputMessageType(SeverityInfo, "Top", b, "percent")
	})
	testing2.AssertLogText(t, "[" + SeverityInfo + "] Some msg 10%---20percent [Top]", output2)
}

func TestOutputMsgType(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputMessageType(SeverityWarning, "Topic", "Msg")
	})
	testing2.AssertLogText(t, fmt.Sprintf("[%s] Msg [Topic]", SeverityWarning), output)
}

func TestOutputMsgWithoutTopic(t *testing.T) {
	output := testing2.CaptureOutput(func() {
		OutputMessageType(SeverityWarning, "", "Msg")
	})
	testing2.AssertLogText(t, fmt.Sprintf("[%s] Msg", SeverityWarning), output)
}

func TestCutMessageWithNoLimit(t *testing.T) {
	SetMaxMessageLength(0)
	msg := CutMessageIfNeeded("some msg")
	assert.Equal(t, "some msg", msg)
}
