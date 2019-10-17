package types

import (
	"github.com/stretchr/testify/assert"
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

	notMatchedFlags := flagSet.GetNotMatchedKeys()
	assert.Equal(t, []string{"red", "green"}, notMatchedFlags)
}
