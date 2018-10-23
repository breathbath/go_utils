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
	for enumValue, mappedValue := range mappedValues {
		if mappedValue == valueToConvert {
			return enumValue, nil
		}
	}

	return 0, fmt.Errorf("Unknown value '%s' for enum '%s'", valueToConvert, enumName)
}

func ConvertInterfaceToEnum(mappedValues map[int]string, valueToConvert interface{}, enumName string) (int, error) {
	if valueToConvert == nil {
		return 0, fmt.Errorf("Empty value to convert for enum '%s'", enumName)
	}
	switch valueToConvert.(type) {
	case string:
		return ConvertStringToEnum(mappedValues, valueToConvert.(string), enumName)
	case []byte:
		valueToConvertString := string(valueToConvert.([]byte))
		return ConvertStringToEnum(mappedValues, valueToConvertString, enumName)
	default:
		return 0, fmt.Errorf("Non-string value '%v' for enum '%s'", valueToConvert, enumName)
	}
}

func GenerateEnumQueryPart(mappedValues map[int]string) string {
	enumQueryValues := []string{}
	for _, val := range mappedValues {
		enumQueryValues = append(enumQueryValues, fmt.Sprintf("'%s'", val))
	}

	return strings.Join(enumQueryValues, ",")
}
