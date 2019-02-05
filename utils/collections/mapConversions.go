//Package collections provides methods to operate on collections of values (structs/maps/arrays)
package collections

import (
	"fmt"
	"strings"
	"sync"
)

//ConvertStructToMap converts []string{"red", "blue", "red"} to map[string]bool{"red":true,"blue":true}
func ConvertStructToMap(input []string) map[string]bool {
	output := map[string]bool{}
	for _, inputItem := range input {
		output[inputItem] = true
	}

	return output
}

//ConvertStructToSyncMap converts []string{"red", "blue"} to sync.Map{"red":true,"blue":true}
func ConvertStructToSyncMap(input []string) *sync.Map {
	output := sync.Map{}
	for _, inputItem := range input {
		output.Store(inputItem, true)
	}

	return &output
}

/*
 * ConvertStructToSyncMapWithCallback does the same as ConvertStructToSyncMap with additional filter/mutator callback
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

//GetMapValueOrError returns error if map doesn't contain the provided key
func GetMapValueOrError(input map[string]string, key string) (string, error) {
	v, ok := input[key]
	if !ok {
		return "", fmt.Errorf("No value for key '%s' in the map", key)
	}

	return v, nil
}

//ExtractMapValues converts map[string]interface{}{"Bob":"Bob", "Alice":"Alice", "John":"John"}
//by filter values []interface{}{"Bob","Alice"} to []interface{}{"Bob", "Alice", "John"}
func ExtractMapValues(inputMap map[string]interface{}) []interface{} {
	result := []interface{}{}
	for _, val := range inputMap {
		result = append(result, val)
	}

	return result
}

//JoinMap converts map[string]string{"name":"Bob","color":"red","size","big"} to
//2 strings "name,color,size" and "Bob,red,big" if "," is provided as separator
func JoinMap(inputMap map[string]string, sep string) (keysStr, valuesStr string) {
	keys, values := "", ""
	for key, val := range inputMap {
		keys += key + sep
		values += val + sep
	}
	return strings.TrimRight(keys, sep), strings.TrimRight(values, sep)
}

//MapToSlice converts map[string]string{"name":"Bob","color":"red","size","big"} to
//keys and values slices like []string{"name","color","size"} and []string{"Bob","red","big"}
func MapToSlices(inputMap map[string]string) (keys, values []string) {
	for key, val := range inputMap {
		keys = append(keys, key)
		values = append(values, val)
	}

	return
}
