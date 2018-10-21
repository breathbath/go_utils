package conv

import "fmt"

func StringerToJson(input fmt.Stringer) []byte {
	jsonValue := fmt.Sprintf(`"%s"`, input.String())
	return []byte(jsonValue)
}
