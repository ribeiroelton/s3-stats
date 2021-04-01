package comparedate

import "time"

func GetMostRecentDate(cur, new *time.Time) *time.Time {
	if new.After(*cur) {
		return new
	}
	return cur
}
