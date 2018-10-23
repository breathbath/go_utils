package time

import "time"

func GetUnixTimestampMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTimeFromTimestampMilliseconds(millisecondsTimeStamp int64) time.Time {
	return time.Unix(millisecondsTimeStamp/1000, 0).UTC()
}
