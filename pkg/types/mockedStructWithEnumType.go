package types

type Person struct {
	Salutation Salutation `json:"salutation"`
	Name       string     `json:"name"`
}
