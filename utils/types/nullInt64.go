package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

type NullInt64 struct {
	sql.NullInt64
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) UnmarshalJSON(input []byte) error {
	var targetInt int64
	var targetStr string

	convErrInt := json.Unmarshal(input, &targetInt)
	convErrStr := json.Unmarshal(input, &targetStr)

	if convErrInt == nil && convErrStr == nil && targetStr == "" && targetInt == 0 {
		ni.Int64 = targetInt
		ni.Valid = false
		return nil
	}

	if convErrInt == nil {
		ni.Int64 = targetInt
		ni.Valid = true
		return nil
	}

	if convErrStr != nil {
		return fmt.Errorf("Cannot convert '%s' to a valid int value", string(input))
	}

	iVal, convErrInt := strconv.ParseInt(targetStr, 10, 64)
	if convErrInt != nil {
		return fmt.Errorf("Cannot convert '%s' to a valid int value", targetStr)
	}

	ni.Int64 = iVal
	ni.Valid = true

	return nil
}
