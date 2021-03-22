package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type StringSet []string

func (vs *StringSet) Add(value string) {
	*vs = append(*vs, value)
}

func (vs *StringSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(vs.ToStrings())
}

func (vs StringSet) ToStrings() (result []string) {
	return []string(vs)
}

func (vs StringSet) Contains(needle string) bool {
	for _, item := range vs {
		if item == needle {
			return true
		}
	}
	return false
}

func (vs StringSet) String() string {
	return strings.Join(vs, ",")
}

func (vs *StringSet) UnmarshalJSON(jsonInput []byte) error {
	var items []string
	err := json.Unmarshal(jsonInput, &items)
	*vs = items
	return err
}

func (vs *StringSet) Scan(value interface{}) error {
	switch strVal := value.(type) {
	case string:
		vs.setFromString(strVal)
	case []byte:
		vs.setFromString(string(strVal))
	default:
		return fmt.Errorf("unknown value type for ValueSet: %v", value)
	}

	return nil
}

func (vs *StringSet) setFromString(rawString string) {
	if rawString == "" {
		return
	}

	*vs = strings.Split(rawString, ",")
}

func (vs StringSet) Value() (driver.Value, error) {
	return vs.String(), nil
}
