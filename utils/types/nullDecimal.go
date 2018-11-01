package types

import (
	"database/sql/driver"
	"encoding/json"
)

type NullDecimal struct {
	DecimalValue Decimal
	Valid        bool
}

func (nd NullDecimal) MarshalJSON() ([]byte, error) {
	if !nd.Valid {
		return json.Marshal(nil)
	}
	return nd.DecimalValue.MarshalJSON()
}

func (nd *NullDecimal) Scan(value interface{}) (err error) {
	isValid := true

	switch value.(type) {
	case nil:
		isValid = false
	case []byte:
		isValid = string(value.([]byte)) != ""
	case string:
		isValid = value.(string) != ""
	}

	if !isValid {
		nd.DecimalValue, nd.Valid = ZERO, false
		return nil
	}

	nd.Valid = isValid

	return nd.DecimalValue.Scan(value)
}

func (nd NullDecimal) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}
	return nd.DecimalValue, nil
}

func (nd *NullDecimal) UnmarshalJSON(input []byte) error {
	d := Decimal{}
	err := d.UnmarshalJSON(input)

	nd.DecimalValue = d
	//will be not valid for all null values, all invalid decimal values should be considered as errors
	nd.Valid = string(input) != "null"

	return err
}
