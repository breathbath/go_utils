package time

import "time"

func GetUnixTimestampMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTimeFromTimestampMilliseconds(millisecondsTimeStamp int64) time.Time {
	milliSecondsTimespamp := time.Millisecond * time.Duration(millisecondsTimeStamp)

	return time.Unix(int64(milliSecondsTimespamp.Seconds()), 0).UTC()
}
