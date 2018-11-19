package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Int64OutputExpectation struct {
	InputStr        string
	ExpectedResult  int64
	ExpectedIsValid bool
	ExpectedError   string
}

type Int64InputExpectation struct {
	InputInt       int64
	InputIsValid   bool
	ExpectedResult string
}

func TestNullInt64JsonOutputConversion(t *testing.T) {
	dataSets := []Int64OutputExpectation{
		{
			InputStr:        `null`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "",
		},
		{
			InputStr:        `2`,
			ExpectedResult:  2,
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `"2"`,
			ExpectedResult:  2,
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `1.2`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "Cannot convert '1.2' to a valid int value",
		},
		{
			InputStr:        `0`,
			ExpectedResult:  0,
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `"0"`,
			ExpectedResult:  0,
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `""`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "Cannot convert '' to a valid int value",
		},
		{
			InputStr:        `" "`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "Cannot convert ' ' to a valid int value",
		},
		{
			InputStr:        `"dfjaslfa"`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "Cannot convert 'dfjaslfa' to a valid int value",
		},
	}

	for k, dataSet := range dataSets {
		errorMsg := fmt.Sprintf("Check %v, test case %d", dataSet, k)
		var val NullInt64
		err := json.Unmarshal([]byte(dataSet.InputStr), &val)
		assert.Equal(t, dataSet.ExpectedResult, val.Int64, errorMsg)
		assert.Equal(t, dataSet.ExpectedIsValid, val.Valid, errorMsg)
		if dataSet.ExpectedError != "" || err != nil {
			assert.EqualError(t, err, dataSet.ExpectedError, errorMsg)
		}
	}
}

func TestNullInt64ToJsonConversion(t *testing.T) {
	dataSets := []Int64InputExpectation{
		{
			InputInt:       0,
			InputIsValid:   true,
			ExpectedResult: "0",
		},
		{
			InputInt:       1,
			InputIsValid:   true,
			ExpectedResult: "1",
		},
		{
			InputInt:       -1,
			InputIsValid:   false,
			ExpectedResult: "null",
		},
		{
			InputInt:       -1,
			InputIsValid:   true,
			ExpectedResult: "-1",
		},
	}

	for k, dataSet := range dataSets {
		errorMsg := fmt.Sprintf("Check %v, test case %d", dataSet, k)

		val := NullInt64{sql.NullInt64{Int64: dataSet.InputInt, Valid: dataSet.InputIsValid}}
		jsonResult, err := json.Marshal(val)
		assert.NoError(t, err, errorMsg)
		assert.Equal(t, dataSet.ExpectedResult, string(jsonResult), errorMsg)
	}
}
