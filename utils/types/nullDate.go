package types

import (
	"time"
)

type NullDate struct {
	NullTime
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
