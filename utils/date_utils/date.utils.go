package date_utils

import (
	"github.com/ocassio/timetable-go-api/models"
	"time"
)

const DateFormat = "02.01.2006"
const TimeFormat = "15:04"

const Day = 24 * time.Hour

var daysOfWeek = []string {
	"Воскресенье",
	"Понедельник",
	"Вторник",
	"Среда",
	"Четверг",
	"Пятница",
	"Суббота",
}

func ToDateString(time time.Time) string {
	return time.Format(DateFormat)
}

func ToTimeString(time time.Time) string {
	return time.Format(TimeFormat)
}

func ToDate(value string) (time.Time, error) {
	return time.Parse(DateFormat, value)
}

func GetDayOfWeekName(time time.Time) string {
	return daysOfWeek[time.Weekday()]
}

func GetSevenDays(f *time.Time) models.DateRange {
	from := f
	if from == nil {
		now := time.Now()
		from = &now
	}

	to := from.Add(6 * Day)

	return models.DateRange {
		From: *from,
		To: to,
	}
}

func GetCurrentWeek() models.DateRange {
	from := GetFirstDayOfWeek(time.Now())
	return GetSevenDays(&from)
}

func GetFirstDayOfWeek(date time.Time) time.Time {
	if date.Weekday() == time.Sunday {
		return date.Add(-6 * Day)
	}

	delta := date.Weekday() - time.Monday
	return date.Add(-time.Duration(delta) * Day)
}
