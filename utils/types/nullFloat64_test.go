package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Float64OutputExpectation struct {
	InputStr        string
	ExpectedResult  float64
	ExpectedIsValid bool
	ExpectedError   string
}

type Float64InputExpectation struct {
	InputFloat     float64
	InputIsValid   bool
	ExpectedResult string
}

func TestNullFloat64ToJsonConversion(t *testing.T) {
	dataSets := []Float64InputExpectation{
		{
			InputFloat:     0,
			InputIsValid:   true,
			ExpectedResult: "0",
		},
		{
			InputFloat:     1.2,
			InputIsValid:   true,
			ExpectedResult: "1.2",
		},
		{
			InputFloat:     1.2,
			InputIsValid:   false,
			ExpectedResult: "null",
		},
		{
			InputFloat:     -12,
			InputIsValid:   true,
			ExpectedResult: "-12",
		},
	}

	for k, dataSet := range dataSets {
		errorMsg := fmt.Sprintf("Check %v, test case %d", dataSet, k)

		val := NullFloat64{sql.NullFloat64{Float64: dataSet.InputFloat, Valid: dataSet.InputIsValid}}
		jsonResult, err := json.Marshal(val)
		assert.NoError(t, err, errorMsg)
		assert.Equal(t, dataSet.ExpectedResult, string(jsonResult), errorMsg)
	}
}

func TestJsonToNullFloat64Conversion(t *testing.T) {
	dataSets := []Float64OutputExpectation{
		{
			InputStr:        `null`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "",
		},
		{
			InputStr:        `2.64`,
			ExpectedResult:  2.64,
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `"2.64"`,
			ExpectedResult:  2.64,
			ExpectedIsValid: true,
			ExpectedError:   "",
		},
		{
			InputStr:        `1`,
			ExpectedResult:  1.0,
			ExpectedIsValid: true,
			ExpectedError:   "",
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
			ExpectedError:   "cannot convert '' to a valid float value",
		},
		{
			InputStr:        `" "`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "cannot convert ' ' to a valid float value",
		},
		{
			InputStr:        `"dafjldfja"`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "cannot convert 'dafjldfja' to a valid float value",
		},
		{
			InputStr:        `[1,2]`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "cannot convert '[1,2]' to a valid float value",
		},
		{
			InputStr:        `{"a":1}`,
			ExpectedResult:  0,
			ExpectedIsValid: false,
			ExpectedError:   "cannot convert '{\"a\":1}' to a valid float value",
		},
	}

	for k, dataSet := range dataSets {
		errorMsg := fmt.Sprintf("Check %v, test case %d", dataSet, k)

		var val NullFloat64
		err := json.Unmarshal([]byte(dataSet.InputStr), &val)
		assert.Equal(t, dataSet.ExpectedResult, val.Float64, errorMsg)
		assert.Equal(t, dataSet.ExpectedIsValid, val.Valid, errorMsg)
		if dataSet.ExpectedError != "" || err != nil {
			assert.EqualError(t, err, dataSet.ExpectedError, errorMsg)
		}
	}
}
