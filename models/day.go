package models

import "strings"

type Day struct {
	Date string `json:"date"`
	DayOfWeek string `json:"dayOfWeek"`
	Lessons []Lesson `json:"lessons"`
}

func (d *Day) Contains(query string) bool {
	lQuery := strings.ToLower(query)
	if contains(d.Date, lQuery) {
		return true
	}

	for _, lesson := range d.Lessons {
		if lesson.Contains(lQuery) {
			return true
		}
	}

	return false
}
