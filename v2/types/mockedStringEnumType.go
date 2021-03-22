package types

import (
	"database/sql/driver"

	"github.com/breathbath/go_utils/v2/conv"
)

type Salutation int

const (
	MS Salutation = iota
	MRS
	DR
	MR
)

var SalutationMap map[int]string

func init() {
	SalutationMap = map[int]string{
		int(MS):  "Miss",
		int(MRS): "Missus",
		int(DR):  "Doctor",
		int(MR):  "Mister",
	}
}

func (s Salutation) MarshalJSON() ([]byte, error) {
	return conv.StringerToJSON(s), nil
}

func (s Salutation) String() string {
	return ConvertEnumToString(SalutationMap, int(s))
}

func (s *Salutation) UnmarshalJSON(jsonInput []byte) error {
	uidName := string(jsonInput)
	return s.SetFromName(uidName)
}

func (s *Salutation) SetFromName(uidName string) error {
	enumVal, err := ConvertStringToEnum(SalutationMap, uidName, "SalutationMap")
	if err != nil {
		return err
	}

	*s = Salutation(enumVal)

	return nil
}

func (s *Salutation) Scan(value interface{}) error {
	enumVal, err := ConvertInterfaceToEnum(SalutationMap, value, "SalutationMap")
	if err != nil {
		return err
	}
	*s = Salutation(enumVal)

	return nil
}

func (s Salutation) Value() (driver.Value, error) {
	return s.String(), nil
}
