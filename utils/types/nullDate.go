package types

import (
	"encoding/json"
	"fmt"
	"time"
)

const NullableStr = "null"
const EmptyStr = `""`

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
	if string(input) == NullableStr || string(input) == EmptyStr {
		nd.Valid = false
		return nil
	}

	var err error
	var t time.Time
	t, err = time.Parse(`"2006-01-02"`, string(input))

	nd.Time = t
	nd.Valid = string(input) != NullableStr

	return err
}
