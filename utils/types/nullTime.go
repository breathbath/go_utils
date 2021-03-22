package types

import (
	"encoding/json"
	"time"

	"github.com/go-sql-driver/mysql"
)

type NullTime struct {
	mysql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(input []byte) error {
	t := time.Time{}
	if string(input) == NullableStr || string(input) == EmptyStr {
		nt.Valid = false
		return nil
	}

	err := t.UnmarshalJSON(input)

	nt.Time = t
	// will be not valid for all null dates, all invalid date values should be considered as errors
	nt.Valid = string(input) != "null"

	return err
}
