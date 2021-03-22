package conv

import (
	"fmt"
	"strings"
)

// This example converts struct to sync map with additional applying of callback function to modify
// and filter elements
func ExampleConvertStructToSyncMapWithCallback() {
	// Want to have unique map of file names without extension and without '.' and '..'
	// something like sync.Map{"file":true}
	collectedFileNames := []string{".", "..", "file.txt", "file.txt"}
	resultMap := ConvertStructToSyncMapWithCallback(
		collectedFileNames,
		func(inputItem string) (modifiedString string, shouldBeIncluded bool) {
			if inputItem == "." || inputItem == ".." {
				return "", false
			}
			return strings.Split(inputItem, ".")[0], true
		},
	)

	printableMap := map[string]bool{}
	resultMap.Range(func(key, value interface{}) bool {
		printableMap[key.(string)] = value.(bool)
		return true
	})

	fmt.Printf("%+v", printableMap)

	// Output:
	// map[file:true]
}
