package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type NullDate struct {
	NullTime
}

func (nd NullDate) MarshalJSON() ([]byte, error) {
	if !nd.Valid {
		return json.Marshal(nil)
	}

	outputStr := fmt.Sprintf(`"%s"`, nd.Time.Format("2006-01-02"))

	return []byte(outputStr), nil
}

func (nd *NullDate) UnmarshalJSON(input []byte) error {
	t := time.Time{}
	if string(input) == "null" {
		return nil
	}

	var err error
	t, err = time.Parse(`"2006-01-02"`, string(input))

	nd.Time = t
	nd.Valid = string(input) != "null"

	return err
}
