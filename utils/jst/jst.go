package jst

import "time"

const (
	// location ...
	location = `Asia/Tokyo`
)

func init() {
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

// Now wraps time.Now
func Now() time.Time {
	return time.Now()
}

// Parse wraps time.ParseInLocation
func Parse(layout string, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}
