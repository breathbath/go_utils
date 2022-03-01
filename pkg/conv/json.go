package conv

import "fmt"

func StringerToJSON(input fmt.Stringer) []byte {
	jsonValue := fmt.Sprintf(`%q`, input.String())
	return []byte(jsonValue)
}
