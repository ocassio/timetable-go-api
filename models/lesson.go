package models

import "strings"

type Lesson struct {
	Number  string    `json:"number"`
	Room    string    `json:"room"`
	Name    string    `json:"name"`
	Teacher string    `json:"teacher"`
	Type    string    `json:"type"`
	Group   string    `json:"group"`
	Note    string    `json:"note"`
	Time    TimeRange `json:"time"`
}

func NewLesson(params *[]string) *Lesson {
	p := *params
	return &Lesson{
		Room:    p[0],
		Number:  p[1],
		Teacher: p[2],
		Type:    p[3],
		Name:    p[4],
		Group:   p[5],
		Note:    p[6],
	}
}

func (l *Lesson) Contains(query string) bool {
	lQuery := strings.ToLower(query)
	return contains(l.Room, lQuery) ||
		contains(l.Name, lQuery) ||
		contains(l.Teacher, lQuery) ||
		contains(l.Type, lQuery) ||
		contains(l.Group, lQuery) ||
		contains(l.Note, lQuery)
}

func contains(target string, query string) bool {
	lowerTarget := strings.ToLower(target)
	return strings.Contains(lowerTarget, query)
}
