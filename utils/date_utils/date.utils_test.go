package date_utils

import (
	"testing"
	"time"
)

func TestGetSevenDays(t *testing.T) {
	date, _ := time.Parse(DateFormat, "07.06.2020")

	result := GetSevenDays(&date)

	expectedTo, _ := time.Parse(DateFormat, "13.06.2020")
	if date != result.From {
		t.Errorf("Expected: %s; Actual: %s", date, result.From)
	}
	if expectedTo != result.To {
		t.Errorf("Expected: %s; Actual: %s", expectedTo, result.To)
	}
}

func TestGetFirstDayOfWeek(t *testing.T) {
	date, _ := time.Parse(DateFormat, "07.06.2020")

	result := GetFirstDayOfWeek(date)

	expected, _ := time.Parse(DateFormat, "01.06.2020")
	if expected != result {
		t.Errorf("Expected: %s; Actual: %s", expected, result)
	}

	date, _ = time.Parse(DateFormat, "08.06.2020")

	result = GetFirstDayOfWeek(date)

	expected, _ = time.Parse(DateFormat, "08.06.2020")
	if expected != result {
		t.Errorf("Expected: %s; Actual: %s", expected, result)
	}

	date, _ = time.Parse(DateFormat, "12.06.2020")

	result = GetFirstDayOfWeek(date)

	expected, _ = time.Parse(DateFormat, "08.06.2020")
	if expected != result {
		t.Errorf("Expected: %s; Actual: %s", expected, result)
	}
}
