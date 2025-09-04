package timex

import "time"

func StartOfDay(t time.Time) time.Time {
	return StartOfDayByTz(t, t.Location())
}

func StartOfDayByTz(t time.Time, loc *time.Location) time.Time {
	year, month, day := t.In(loc).Date()

	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func IsStartOfDay(t time.Time) bool {
	return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0
}

func IsStartOfDayByTz(t time.Time, loc *time.Location) bool {
	t = t.In(loc)
	return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0
}
