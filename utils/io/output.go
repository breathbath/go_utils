package io

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var maxMessageLength int = 5000

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
