package http

import (
	"net/http"
	"strconv"
)

func GetRequestValueString(req *http.Request, key, defaultValue string) string {
	val, ok := req.URL.Query()[key]
	if !ok {
		return defaultValue
	}

	return val[0]
}

func GetRequestValueInt(req *http.Request, key string, defaultValue int64) int64 {
	val, ok := req.URL.Query()[key]
	if !ok {
		return defaultValue
	}

	valInt, err := strconv.Atoi(val[0])
	if err == nil {
		return int64(valInt)
	}

	return defaultValue
}
