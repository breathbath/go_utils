package conv

import "fmt"

func StringerToJSON(input fmt.Stringer) []byte {
	jsonValue := fmt.Sprintf(`"%s"`, input.String())
	return []byte(jsonValue)
}
