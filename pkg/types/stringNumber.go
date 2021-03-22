package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StringNumber struct {
	numb string
}

func (sn StringNumber) MarshalJSON() ([]byte, error) {
	if sn.numb == "" {
		return json.Marshal(nil)
	}
	return []byte(sn.numb), nil
}

func (sn *StringNumber) UnmarshalJSON(input []byte) error {
	tempNumber := string(input)
	if tempNumber == "" {
		sn.numb = ""
		return nil
	}

	_, err := strconv.ParseFloat(tempNumber, 64)
	if err == nil {
		sn.numb = tempNumber
		return nil
	}

	return fmt.Errorf("cannot unmarshal '%v' into a number type", tempNumber)
}

func (sn StringNumber) String() string {
	return sn.numb
}
