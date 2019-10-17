package types

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestAllTrueFlags(t *testing.T) {
	flagSet := NewFlagSet()

	assert.False(t, flagSet.AllTrue())

	flagSet["red"] = true
	assert.True(t, flagSet.AllTrue())

	flagSet["red"] = false
	assert.False(t, flagSet.AllTrue())
}

func TestGetNotMatchingKeys(t *testing.T) {
	flagSet := FlagSet{
		"red":   false,
		"blue":  true,
		"green": false,
	}

	expectedNotMatchedFlags := []string{"red", "green"}
	notMatchedFlags := flagSet.GetNotMatchedKeys()

	sort.Strings(notMatchedFlags)
	sort.Strings(expectedNotMatchedFlags)

	assert.Equal(t, expectedNotMatchedFlags, notMatchedFlags)
}
