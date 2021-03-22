package types

import (
	"fmt"
	"strings"
)

func ConvertEnumToString(mappedValues map[int]string, valueToConvert int) string {
	strVal, ok := mappedValues[valueToConvert]

	if !ok {
		return ""
	}

	return strVal
}

func ConvertStringToEnum(mappedValues map[int]string, valueToConvert, enumName string) (int, error) {
	valueToConvert = strings.Trim(valueToConvert, `"`)
	for enumValue, mappedValue := range mappedValues {
		if mappedValue == valueToConvert {
			return enumValue, nil
		}
	}

	return 0, fmt.Errorf("unknown value '%s' for enum '%s'", valueToConvert, enumName)
}

func ConvertInterfaceToEnum(mappedValues map[int]string, valueToConvert interface{}, enumName string) (int, error) {
	if valueToConvert == nil {
		return 0, fmt.Errorf("empty value to convert for enum '%s'", enumName)
	}
	switch val := valueToConvert.(type) {
	case string:
		return ConvertStringToEnum(mappedValues, val, enumName)
	case []byte:
		valueToConvertString := string(val)
		return ConvertStringToEnum(mappedValues, valueToConvertString, enumName)
	default:
		return 0, fmt.Errorf("non-string value '%v' for enum '%s'", valueToConvert, enumName)
	}
}

func GenerateEnumQueryPart(mappedValues map[int]string) string {
	enumQueryValues := []string{}
	for _, val := range mappedValues {
		enumQueryValues = append(enumQueryValues, fmt.Sprintf("'%s'", val))
	}

	return strings.Join(enumQueryValues, ",")
}
