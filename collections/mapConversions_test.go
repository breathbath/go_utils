package collections

import (
	testing2 "github.com/breathbath/go_utils/testing"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var exampleStruct []string
var exampleMap map[string]string

func init() {
	exampleStruct = []string{
		"one",
		"two",
		"two",
		"three",
		"one",
	}

	exampleMap = map[string]string{
		"one":   "1",
		"two":   "2",
	}
}

func TestConvertStructToMap(t *testing.T) {
	outputMap := ConvertStructToMap(exampleStruct)
	expectedMap := map[string]bool{
		"one":   true,
		"two":   true,
		"three": true,
	}

	assert.Equal(t, expectedMap, outputMap)
}

func TestConvertStructToSyncMap(t *testing.T) {
	outputMap := ConvertStructToSyncMap(exampleStruct)
	expectedMap := sync.Map{}
	for _, structItem := range exampleStruct {
		expectedMap.Store(structItem, true)
	}

	testing2.AssertSyncMapEqual(t, &expectedMap, outputMap)
}

func TestConvertStructToSyncMapWithCallback(t *testing.T) {
	mapFunc := func(inputItem string) (string, bool) {
		if inputItem == "one" {
			return "", false
		}

		return inputItem + "1", true
	}
	outputMap := ConvertStructToSyncMapWithCallback(exampleStruct, mapFunc)

	expectedMap := sync.Map{}
	expectedMap.Store("two1", true)
	expectedMap.Store("three1", true)

	testing2.AssertSyncMapEqual(t, &expectedMap, outputMap)
}

func TestGetMapValueOrError(t *testing.T) {
	val, err := GetMapValueOrError(exampleMap, "one")
	assert.NoError(t, err)
	assert.Equal(t, "1", val)

	val, err = GetMapValueOrError(exampleMap, "nonExistingKey")
	assert.EqualError(t, err, "No value for key 'nonExistingKey' in the map")
	assert.Equal(t, "", val)
}

func TestExtractMapValues(t *testing.T) {
	filter := []string{"one"}
	actualResult, err := ExtractMapValues(filter, exampleMap)
	assert.NoError(t, err)
	expectedResult := []string{"1"}
	assert.Equal(t, expectedResult, actualResult)

	filter = []string{"two", "nonExistingKey"}
	actualResult, err = ExtractMapValues(filter, exampleMap)
	assert.EqualError(t, err, "No value for key 'nonExistingKey' in the map")
	expectedResult = []string{"2", ""}
	assert.Equal(t, expectedResult, actualResult)
}

func TestJoinMap(t *testing.T) {
	keysStr, valueStr := JoinMap(exampleMap, ",")
	possibleVariantsKeys := map[string] bool {
		"one,two": true,
		"two,one": true,
	}
	possibleVariantsValues := map[string] bool {
		"1,2": true,
		"2,1": true,
	}

	_, ok := possibleVariantsKeys[keysStr]
	if !ok {
		t.Errorf("Expected result should be either 'one,two', or 'two,one' but %s is returned", keysStr)
	}

	_, ok = possibleVariantsValues[valueStr]
	if !ok {
		t.Errorf("Expected result should be either '1,2', or '2,1' but %s is returned", valueStr)
	}
}

func TestMapToStruct(t *testing.T) {
	actualKeys, actualValues := MapToStruct(exampleMap)
	//sort.Strings(actualKeys)
	//sort.Strings(actualValues)

	expectedKeys := []string{"one", "two"}
	expectedValues := []string{"1", "2"}

	assert.ElementsMatch(t, expectedKeys, actualKeys)
	assert.ElementsMatch(t, expectedValues, actualValues)
}
