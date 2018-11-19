package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(input []byte) error {
	var targetFloat float64
	var targetStr string

	convErrFloat := json.Unmarshal(input, &targetFloat)
	convErrStr := json.Unmarshal(input, &targetStr)

	if convErrFloat == nil && convErrStr == nil && targetStr == "" && targetFloat == 0 {
		ns.String = ""
		ns.Valid = false
		return nil
	}

	if convErrFloat == nil {
		ns.String = fmt.Sprint(targetFloat)
		ns.Valid = true
		return nil
	}

	if convErrStr != nil {
		return fmt.Errorf("Cannot convert '%s' to a valid string value", string(input))
	}

	ns.String = targetStr
	ns.Valid = targetStr != ""

	return nil
}
