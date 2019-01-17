package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringSetAdding(t *testing.T) {
	ss := &StringSet{}
	ss.Add("one")
	ss.Add("two")

	assert.EqualValues(t, "one,two", ss.String())
}

func TestStringSetMarshalJson(t *testing.T) {
	ss := &StringSet{"one", "two"}
	jsonSs, err := ss.MarshalJSON()

	assert.NoError(t, err)

	assert.EqualValues(t, `["one","two"]`, jsonSs)
}

func TestStringSetToStringsConversion(t *testing.T) {
	ss := &StringSet{"one", "two"}

	assert.EqualValues(t, []string{"one", "two"}, ss.ToStrings())
}

func TestStringSetUnMarshalJson(t *testing.T) {
	ssJson := []byte(`["one","two"]`)

	providedSs := &StringSet{}
	err := providedSs.UnmarshalJSON(ssJson)
	assert.NoError(t, err)

	expectedSs := &StringSet{"one", "two"}

	assert.EqualValues(t, expectedSs, providedSs)
}

func TestStringSetScan(t *testing.T) {
	providedSs := &StringSet{}

	err := providedSs.Scan("")
	assert.NoError(t, err)
	assert.EqualValues(t, &StringSet{}, providedSs)

	err = providedSs.Scan("one,two")
	assert.NoError(t, err)
	assert.EqualValues(t, &StringSet{"one", "two"}, providedSs)

	notSplitableInputs := []string{"one.two", "one two", "one-two", "one"}
	for _, notSplitableInput := range notSplitableInputs {
		err = providedSs.Scan(notSplitableInput)
		assert.NoError(t, err)
		assert.EqualValues(t, &StringSet{notSplitableInput}, providedSs)
	}

	err = providedSs.Scan([]byte("one,two"))
	assert.NoError(t, err)
	assert.EqualValues(t, &StringSet{"one", "two"}, providedSs)
}

func TestStringSetValueConversion(t *testing.T) {
	ss := &StringSet{"one", "two"}
	val, err := ss.Value()
	assert.NoError(t, err)

	assert.Equal(t, "one,two", val)
}

func TestStringSetSearching(t *testing.T) {
	ss := &StringSet{"one", "two"}
	assert.False(t, ss.Contains("three"))
	assert.True(t, ss.Contains("one"))
	assert.True(t, ss.Contains("two"))
}