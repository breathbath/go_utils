package testing

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func AssertLogText(t *testing.T, expectedText, actualText string) {
	resultWithoutTimeStamp := regexp.MustCompile(`(?m)^\d{4}\/\d{2}\/\d{2}\s\d{2}:\d{2}:\d{2}\s*`).ReplaceAllString(actualText, "")
	resultWithoutTimeStamp = strings.Replace(resultWithoutTimeStamp, "\n", "", -1)
	resultWithoutTimeStamp = strings.TrimSpace(resultWithoutTimeStamp)
	assert.Equal(t, expectedText, resultWithoutTimeStamp)
}
