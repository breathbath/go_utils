package conv

import (
	"sort"
	"sync"
	"testing"

	testing2 "github.com/breathbath/go_utils/utils/testing"
	"github.com/stretchr/testify/assert"
)

var exampleStruct []string
var exampleMap map[string]string
var exampleMapInterface map[string]interface{}

type syncMapTestCase struct {
	syncMap   *sync.Map
	legacyMap map[string]interface{}
	strMap    string
	errStr    string
}

func init() {
	exampleStruct = []string{
		"one",
		"two",
		"two",
		"three",
		"one",
	}

	exampleMap = map[string]string{
		"one": "1",
		"two": "2",
	}

	exampleMapInterface = map[string]interface{}{
		"1": 1,
		"2": 2,
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
	assert.EqualError(t, err, "no value for key 'nonExistingKey' in the map")
	assert.Equal(t, "", val)
}

func TestExtractMapValues(t *testing.T) {
	actualResult := ExtractMapValues(exampleMapInterface)
	expectedResult := []int{1, 2}

	actualResultInts := []int{}
	for _, actualResultI := range actualResult {
		actualResultInts = append(actualResultInts, actualResultI.(int))
	}

	sort.Ints(actualResultInts)

	assert.Equal(t, expectedResult, actualResultInts)
}

func TestJoinMap(t *testing.T) {
	keysStr, valueStr := JoinMap(exampleMap, ",")
	possibleVariantsKeys := map[string]bool{
		"one,two": true,
		"two,one": true,
	}
	possibleVariantsValues := map[string]bool{
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

func TestMapToSlices(t *testing.T) {
	actualKeys, actualValues := MapToSlices(exampleMap)
	sort.Strings(actualKeys)
	sort.Strings(actualValues)

	expectedKeys := []string{"one", "two"}
	expectedValues := []string{"1", "2"}

	assert.ElementsMatch(t, expectedKeys, actualKeys)
	assert.ElementsMatch(t, expectedValues, actualValues)
}

func TestConvertSyncMapToMap(t *testing.T) {
	testCases := []syncMapTestCase{
		{
			syncMap: ConvertMapToSyncMap(map[string]interface{}{
				"intVal":  1,
				"boolVal": true,
				"strVal":  "str",
				"nilVal":  nil,
				"mapVal": map[int]int{
					1: 10,
					2: 20,
				},
			}),
			legacyMap: map[string]interface{}{
				"intVal":  1,
				"boolVal": true,
				"strVal":  "str",
				"nilVal":  nil,
				"mapVal": map[int]int{
					1: 10,
					2: 20,
				},
			},
		},
		{
			syncMap:   &sync.Map{},
			legacyMap: map[string]interface{}{},
		},
	}
	for _, testCase := range testCases {
		actualConvertedMap := ConvertSyncMapToMap(testCase.syncMap)
		assert.Equal(t, testCase.legacyMap, actualConvertedMap)
	}
}

func TestConvertMapToSyncMap(t *testing.T) {
	testCases := []syncMapTestCase{
		{
			syncMap: ConvertMapToSyncMap(map[string]interface{}{
				"intVal":  1,
				"boolVal": true,
				"strVal":  "str",
				"nilVal":  nil,
				"mapVal": map[int]int{
					1: 10,
					2: 20,
				},
			}),
			legacyMap: map[string]interface{}{
				"intVal":  1,
				"boolVal": true,
				"strVal":  "str",
				"nilVal":  nil,
				"mapVal": map[int]int{
					1: 10,
					2: 20,
				},
			},
		},
		{
			syncMap:   &sync.Map{},
			legacyMap: map[string]interface{}{},
		},
	}
	for _, testCase := range testCases {
		actualSyncMap := ConvertMapToSyncMap(testCase.legacyMap)
		actualSyncMapConvertedToSimpleMap := ConvertSyncMapToMap(actualSyncMap)
		assert.Equal(t, testCase.legacyMap, actualSyncMapConvertedToSimpleMap)
	}
}

func TestMarshalSyncMap(t *testing.T) {
	testCases := []syncMapTestCase{
		{
			syncMap: ConvertMapToSyncMap(map[string]interface{}{
				"intVal":  1,
				"boolVal": true,
				"strVal":  "str",
				"nilVal":  nil,
				"mapVal": map[int]int{
					1: 10,
					2: 20,
				},
			}),
			strMap: `{"boolVal":true,"intVal":1,"mapVal":{"1":10,"2":20},"nilVal":null,"strVal":"str"}`,
		},
		{
			syncMap: &sync.Map{},
			strMap:  "{}",
		},
	}

	for _, testCase := range testCases {
		mapBytes, err := MarshalSyncMap(testCase.syncMap)
		assert.NoError(t, err)
		if err != nil {
			continue
		}

		assert.Equal(t, testCase.strMap, string(mapBytes))
	}
}

func TestUnMarshalSyncMap(t *testing.T) {
	testCases := []syncMapTestCase{
		{
			syncMap: ConvertMapToSyncMap(map[string]interface{}{
				"intVal":  1,
				"boolVal": true,
				"strVal":  "str",
				"nilVal":  nil,
				"mapVal": map[int]int{
					1: 10,
					2: 20,
				},
			}),
			strMap: `{"boolVal":true,"intVal":1,"mapVal":{"1":10,"2":20},"nilVal":null,"strVal":"str"}`,
			legacyMap: map[string]interface{}{
				"intVal":  1.00,
				"boolVal": true,
				"strVal":  "str",
				"nilVal":  nil,
				"mapVal": map[string]interface{}{
					"1": 10.00,
					"2": 20.00,
				},
			},
		},
		{
			syncMap:   &sync.Map{},
			strMap:    "{}",
			legacyMap: map[string]interface{}{},
		},
		{
			syncMap: &sync.Map{},
			strMap:  "dfadsfas",
			errStr:  "failed to convert `dfadsfas` to map[string]interface{}",
		},
	}

	for _, testCase := range testCases {
		unmarshalledSyncMap, err := UnMarshalSyncMap([]byte(testCase.strMap))
		if testCase.errStr != "" {
			assert.Error(t, err)
			if err != nil {
				assert.Contains(t, err.Error(), testCase.errStr)
			}
			return
		}

		assert.NoError(t, err)
		if err != nil {
			continue
		}

		unmarshalledSyncMapConvertedToMap := ConvertSyncMapToMap(unmarshalledSyncMap)

		assert.Equal(t, testCase.legacyMap, unmarshalledSyncMapConvertedToMap)
	}
}
