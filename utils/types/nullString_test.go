package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StringOutputExpectation struct {
	InputStr        string
	ExpectedResult  string
	ExpectedIsValid bool
	ExpectedError   string
}

type StringInputExpectation struct {
	InputString    string
	InputIsValid   bool
	ExpectedResult string
}

func TestNullStringJsonOutputConversion(t *testing.T) {
	dataSets := []StringOutputExpectation{
		{
			InputStr:        `null`,
			ExpectedResult:  "",
			ExpectedIsValid: false,
			ExpectedError:   "",
		},
		{
			InputStr:        `2`,
			ExpectedResult:  "2",
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `"2"`,
			ExpectedResult:  "2",
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `1.2`,
			ExpectedResult:  "1.2",
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `0`,
			ExpectedResult:  "0",
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `"0"`,
			ExpectedResult:  "0",
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `""`,
			ExpectedResult:  "",
			ExpectedIsValid: false,
			ExpectedError:   "",
		},
		{
			InputStr:        `" "`,
			ExpectedResult:  " ",
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `[1,2]`,
			ExpectedResult:  "",
			ExpectedIsValid: false,
			ExpectedError:   "cannot convert '[1,2]' to a valid string value",
		},
		{
			InputStr:        `{"a":12}`,
			ExpectedResult:  "",
			ExpectedIsValid: false,
			ExpectedError:   "cannot convert '{\"a\":12}' to a valid string value",
		},
	}

	for k, dataSet := range dataSets {
		errorMsg := fmt.Sprintf("Check %v, test case %d", dataSet, k)
		var val NullString
		err := json.Unmarshal([]byte(dataSet.InputStr), &val)
		assert.Equal(t, dataSet.ExpectedResult, val.String, errorMsg)
		assert.Equal(t, dataSet.ExpectedIsValid, val.Valid, errorMsg)
		if dataSet.ExpectedError != "" || err != nil {
			assert.EqualError(t, err, dataSet.ExpectedError, errorMsg)
		}
	}
}

func TestNullStringToJsonConversion(t *testing.T) {
	dataSets := []StringInputExpectation{
		{
			InputString:    "0",
			InputIsValid:   true,
			ExpectedResult: `"0"`,
		},
		{
			InputString:    "",
			InputIsValid:   true,
			ExpectedResult: `""`,
		},
		{
			InputString:    "someStr",
			InputIsValid:   false,
			ExpectedResult: "null",
		},
		{
			InputString:    "someStr",
			InputIsValid:   true,
			ExpectedResult: `"someStr"`,
		},
	}

	for k, dataSet := range dataSets {
		errorMsg := fmt.Sprintf("Check %v, test case %d", dataSet, k)

		val := NullString{sql.NullString{String: dataSet.InputString, Valid: dataSet.InputIsValid}}
		jsonResult, err := json.Marshal(val)
		assert.NoError(t, err, errorMsg)
		assert.Equal(t, dataSet.ExpectedResult, string(jsonResult), errorMsg)
	}
}
