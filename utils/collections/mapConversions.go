package collections

import (
	"fmt"
	"github.com/breathbath/go_utils/utils/errs"
	"strings"
	"sync"
)

func ConvertStructToMap(input []string) map[string]bool {
	output := map[string]bool{}
	for _, inputItem := range input {
		output[inputItem] = true
	}

	return output
}

func ConvertStructToSyncMap(input []string) *sync.Map {
	output := sync.Map{}
	for _, inputItem := range input {
		output.Store(inputItem, true)
	}

	return &output
}

/*
 * callback result string should return a modified value, result bool if false item won't be included in the output
 */
func ConvertStructToSyncMapWithCallback(input []string, callback func(inputItem string) (string, bool)) *sync.Map {
	output := sync.Map{}
	for _, inputItem := range input {
		formattedItem, shouldInclude := callback(inputItem)
		if !shouldInclude {
			continue
		}
		output.Store(formattedItem, true)
	}

	return &output
}

func GetMapValueOrError(input map[string]string, key string) (string, error) {
	v, ok := input[key]
	if !ok {
		return "", fmt.Errorf("No value for key '%s' in the map", key)
	}

	return v, nil
}

func ExtractMapValues(keys []string, input map[string]string) ([]string, error) {
	result := []string{}
	errsContainer := errs.NewErrorContainer()
	for _, k := range keys {
		v, e := GetMapValueOrError(input, k)
		errsContainer.AddError(e)
		result = append(result, v)
	}

	return result, errsContainer.Result(". ")
}

func JoinMap(inputMap map[string]string, sep string) (keysStr, valuesStr string) {
	keys, values := MapToStruct(inputMap)
	keysStr, valuesStr = strings.Join(keys, sep), strings.Join(values, sep)

	return
}

func MapToStruct(inputMap map[string]string) (keys, values []string) {
	for key, val := range inputMap {
		keys = append(keys, key)
		values = append(values, val)
	}

	return
}
