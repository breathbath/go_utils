package types

type StringNumber struct {
	numb string
}

func (sn StringNumber) MarshalJSON() ([]byte, error) {
	return []byte(sn.numb), nil
}

func (sn *StringNumber) UnmarshalJSON(input []byte) error {
	sn.numb = string(input)
	return nil
}

func (sn StringNumber) String() string {
	return sn.numb
}
