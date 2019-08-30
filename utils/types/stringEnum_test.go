package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertEnumToString(t *testing.T) {
	assert.Equal(t, "Mister", MR.String())
	assert.Equal(t, "Missus", MRS.String())
	assert.Equal(t, "Doctor", DR.String())
	assert.Equal(t, "Miss", MS.String())

	assert.Equal(t, "", ConvertEnumToString(SalutationMap, -1))
}

func TestConvertStructWithEnumToJson(t *testing.T) {
	person := Person{
		Salutation: MR,
		Name: "Dow Jones",
	}

	jsonPerson, err := json.Marshal(person)
	assert.NoError(t, err)
	expectedStringValue := `{"salutation":"Mister","name":"Dow Jones"}`
	actualStringValue := string(jsonPerson)

	assert.Equal(t, expectedStringValue, actualStringValue)
}

func TestConvertStructWithEnumFromJson(t *testing.T) {
	actualPerson := Person{}

	err := json.Unmarshal([]byte(`{"salutation":"Mister","name":"Dow Jones"}`), &actualPerson)
	assert.NoError(t, err)

	expectedPerson := Person{
		Salutation: MR,
		Name: "Dow Jones",
	}

	assert.Equal(t, expectedPerson, actualPerson)
}

func TestConvertStringToEnum(t *testing.T) {
	var someSalutation Salutation
	err := someSalutation.SetFromName("Mister")

	assert.NoError(t, err)
	assert.Equal(t, MR, someSalutation)

	driverVal, err := someSalutation.Value()
	assert.NoError(t, err)
	assert.NotNil(t, driverVal)

	val, err := ConvertStringToEnum(SalutationMap, "Some unknown value", "SalutationMap")
	assert.EqualError(t, err, "Unknown value 'Some unknown value' for enum 'SalutationMap'")
	assert.Equal(t, 0, val)

	err = someSalutation.SetFromName("lala")
	assert.Error(t, err)
}

func TestConvertInterfaceToEnum(t *testing.T) {
	var someSalutation Salutation

	err := someSalutation.Scan("Doctor")
	assert.NoError(t, err)
	assert.Equal(t, DR, someSalutation)

	err = someSalutation.Scan([]byte("Miss"))
	assert.NoError(t, err)
	assert.Equal(t, MS, someSalutation)

	err = someSalutation.Scan(nil)
	assert.EqualError(t, err, "Empty value to convert for enum 'SalutationMap'")

	err = someSalutation.Scan(123)
	assert.EqualError(t, err, "Non-string value '123' for enum 'SalutationMap'")
}

func TestGenerateEnumQueryPart(t *testing.T) {
	salutationEnumSql := GenerateEnumQueryPart(SalutationMap)
	assert.Contains(t, salutationEnumSql, "Doctor")
	assert.Contains(t, salutationEnumSql, "Miss")
	assert.Contains(t, salutationEnumSql, "Missus")
	assert.Contains(t, salutationEnumSql, "Mister")
}
