package io

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var maxMessageLength = 5000

var logger Logger

type Logger interface {
	OutputMessageType(messageType, topic, msg string, args ...interface{})
}

func SetLogger(l Logger) {
	logger = l
}

const (
	SeverityError   = "ERROR"
	SeverityInfo    = "INFO"
	SeverityWarning = "WARNING"
)

type DefaultLogger struct{}

func (dl DefaultLogger) OutputMessageType(messageType, topic, msg string, args ...interface{}) {
	finalMsg := msg
	if len(args) > 0 {
		finalMsg = fmt.Sprintf(msg, args...)
	}
	finalMsg = CutMessageIfNeeded(finalMsg)

	finalMsg = GenerateMessage(
		messageType,
		finalMsg,
		topic,
	)

	finalMsg = RemoveLineBreaks(finalMsg)

	log.Println(finalMsg)
}

/*
OutputError shows error as
[ERROR] {error_message}, {message} [{topic}] e.g.
[ERROR] Cannot write file "tt.txt", please check permissions 664 on folder /tmp/tt.txt [order_123]
topic is needed to group logged event for certain process e.g. all events which happened to order_123
we also cut line separators because normally such messages are split by them as separate messages
we also cut too long messages (more than maxMessageLength chars) because of the processing overhead for long outputs
*/
func OutputError(err error, topic, msg string, args ...interface{}) {
	msgToOutput := fmt.Sprintf(msg, args...)
	finalOutput := ""
	if msgToOutput == "" {
		finalOutput = err.Error()
	} else {
		finalOutput = fmt.Sprintf("%v, %s", err, msgToOutput)
	}

	OutputMessageType(SeverityError, topic, finalOutput)
}

/*
OutputWarning shows warning as
[WARNING] {message} [{topic}] e.g.
[WARNING] Cannot send email [order_123]
*/
func OutputWarning(topic, msg string, args ...interface{}) {
	OutputMessageType(SeverityWarning, topic, msg, args...)
}

/*
OutputWarning shows warning as
[INFO] {message} [{topic}] e.g.
[INFO] Cannot send email [order_123]
*/
func OutputInfo(topic, msg string, args ...interface{}) {
	OutputMessageType(SeverityInfo, topic, msg, args...)
}

// OutputMessageType shows a message as [{messageType}] {message} [{topic}]
func OutputMessageType(messageType, topic, msg string, args ...interface{}) {
	if logger == nil {
		logger = DefaultLogger{}
	}

	logger.OutputMessageType(messageType, topic, msg, args...)
}

func GenerateMessage(eventType, message, topic string) string {
	if topic == "" && eventType == "" {
		return message
	}

	if topic == "" {
		return fmt.Sprintf("[%s] %s", eventType, message)
	}

	return fmt.Sprintf("[%s] %s [%s]", eventType, message, topic)
}

// OutputSingleLine is deprecated should be used in internal outputs
func OutputSingleLine(message string, args ...interface{}) {
	OutputInfo("", message, args...)
}

// OutputSingleLineWithTopic is deprecated should be used in internal outputs
func OutputSingleLineWithTopic(topic, message string, args ...interface{}) {
	OutputInfo(topic, message, args...)
}

func RemoveLineBreaks(input string) string {
	input = strings.Trim(input, `\n "'`)

	re := regexp.MustCompile(`\r?\n`)
	input = re.ReplaceAllString(input, " ")

	return input
}

func SetMaxMessageLength(newValue int) {
	maxMessageLength = newValue
}

func CutMessageIfNeeded(message string) string {
	if maxMessageLength <= 0 {
		return message
	}

	if len(message) > maxMessageLength {
		return message[0:maxMessageLength]
	}

	return message
}
