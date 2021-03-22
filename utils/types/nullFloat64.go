package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nf.Float64)
}

func (nf *NullFloat64) UnmarshalJSON(input []byte) error {
	var targetFloat float64
	var targetStr string

	convErrFloat := json.Unmarshal(input, &targetFloat)
	convErrStr := json.Unmarshal(input, &targetStr)

	if convErrFloat == nil && convErrStr == nil && targetStr == "" && targetFloat == 0 {
		nf.Float64 = targetFloat
		nf.Valid = false
		return nil
	}

	if convErrFloat == nil {
		nf.Float64 = targetFloat
		nf.Valid = true
		return nil
	}

	if convErrStr != nil {
		return fmt.Errorf("cannot convert '%s' to a valid float value", string(input))
	}

	fVal, convErrFloat := strconv.ParseFloat(targetStr, 64)
	if convErrFloat != nil {
		return fmt.Errorf("cannot convert '%s' to a valid float value", targetStr)
	}

	nf.Float64 = fVal
	nf.Valid = true

	return nil
}
