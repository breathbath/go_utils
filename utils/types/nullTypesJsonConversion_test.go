package types

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockUser struct {
	BirthDay NullTime    `json:"BirthDay"`
	Age      NullDecimal `json:"Age"`
	Name     string      `json:"Name"`
}

type MockUserTest struct {
	jsonInput     string
	expectedError string
	expectedUser  MockUser
}

func TestMockUserFromJsonConversion(t *testing.T) {
	testItems := []MockUserTest{
		{
			jsonInput:     `{"BirthDay":"2001-01-01T11:00:00Z","Age":22,"Name":"Tom"}`,
			expectedError: "",
			expectedUser: MockUser{
				BirthDay: NullTime{
					NullTime: mysql.NullTime{
						Time:  time.Date(2001, 1, 1, 11, 0, 0, 0, time.UTC),
						Valid: true,
					},
				},
				Age: NullDecimal{
					DecimalValue: NewDecimalFromString("22"),
					Valid:        true,
				},
				Name: "Tom",
			},
		},
		{
			jsonInput:     `{"Name":"Jim"}`,
			expectedError: "",
			expectedUser: MockUser{
				BirthDay: NullTime{
					NullTime: mysql.NullTime{
						Time:  time.Time{},
						Valid: false,
					},
				},
				Age: NullDecimal{
					DecimalValue: Decimal{},
					Valid:        false,
				},
				Name: "Jim",
			},
		},
		{
			jsonInput:     `{"BirthDay":null,"Age":null,"Name":"bruno"}`,
			expectedError: "",
			expectedUser: MockUser{
				BirthDay: NullTime{
					NullTime: mysql.NullTime{
						Time:  time.Time{},
						Valid: false,
					},
				},
				Age: NullDecimal{
					DecimalValue: Decimal{},
					Valid:        false,
				},
				Name: "bruno",
			},
		},
	}

	for _, testItem := range testItems {
		var actualUser MockUser
		err := json.Unmarshal([]byte(testItem.jsonInput), &actualUser)
		if testItem.expectedError == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, testItem.expectedError)
		}

		assert.Equal(t, testItem.expectedUser, actualUser)
	}
}
