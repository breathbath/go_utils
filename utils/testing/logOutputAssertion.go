package testing

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func AssertLogText(t *testing.T, expectedText, actualText string) {
	regx := fmt.Sprintf(`\d{4}\/\d{2}\/\d{2}\s\d{2}:\d{2}:\d{2}\s%s`, regexp.QuoteMeta(expectedText))
	assert.Regexp(t, regx, actualText)
}
