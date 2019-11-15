package timingwheel

import "time"

func timeToMS(t time.Time) int64 {

	return t.UnixNano() / int64(time.Millisecond)
}

func truncate(x int64, m int64) int64 {
	if m <= 0 {
		return x
	}
	return x - x%m
}

func msToTime(t int64) time.Time {

	return time.Unix(0, t*int64(time.Millisecond)).UTC()
}
