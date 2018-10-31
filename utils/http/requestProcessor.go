package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const REQUEST_TIME_FORMAT = "2006-01-02T15-04-05"

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

func GetRequestValueTimeWithError(req *http.Request, key string) (time.Time, error) {
	val, ok := req.URL.Query()[key]
	if !ok {
		return time.Time{}, fmt.Errorf("No time value provided for key %s", key)
	}

	timeResult, err := time.Parse(REQUEST_TIME_FORMAT, val[0])
	if err != nil {
		return time.Time{}, err
	}

	return timeResult, nil
}

func GetRequestValueTime(req *http.Request, key string, defaultValue time.Time) time.Time {
	timeVal, err := GetRequestValueTimeWithError(req, key)
	if err != nil {
		return defaultValue
	}

	return timeVal
}
