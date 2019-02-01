package io

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var maxMessageLength int = 5000

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
	msgToOutput = fmt.Sprintf("%v, %s", err, msgToOutput)

	OutputMessageType("ERROR", topic, msgToOutput)
}

/*
OutputWarning shows warning as
[WARNING] {message} [{topic}] e.g.
[WARNING] Cannot send email [order_123]
 */
func OutputWarning(err error, topic, msg string, args ...interface{}) {
	OutputMessageType("WARNING", topic, msg, args...)
}

/*
OutputWarning shows warning as
[INFO] {message} [{topic}] e.g.
[INFO] Cannot send email [order_123]
 */
func OutputInfo(topic, msg string, args ...interface{}) {
	OutputMessageType("INFO", topic, msg, args...)
}

//OutputMessageType shows a message as [{messageType}] {message} [{topic}]
func OutputMessageType(messageType, topic, msg string, args ...interface{}) {
	msgToOutput := GenerateMessage(
		messageType,
		fmt.Sprintf(msg, args...),
		topic,
	)

	OutputSingleLine(msgToOutput)
}

func GenerateMessage(eventType, message, topic string) string {
	return fmt.Sprintf("[%s] %s [%s]", eventType, message, topic)
}

func OutputSingleLine(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	msg = RemoveLineBreaks(msg)
	log.Println(CutMessageIfNeeded(msg))
}

func OutputSingleLineWithTopic(topic, message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	msg = fmt.Sprintf("[%s] %s", topic, msg)

	OutputSingleLine(msg)
}

func RemoveLineBreaks(input string) string {
	input = strings.Trim(input, `\n "'`)

	re := regexp.MustCompile(`\r?\n`)
	input = re.ReplaceAllString(input, " ")

	input = CutMessageIfNeeded(input)

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
