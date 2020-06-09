package data_provider

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/ocassio/timetable-go-api/config"
	"github.com/ocassio/timetable-go-api/models"
	"github.com/ocassio/timetable-go-api/utils/date_utils"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const colCount = 7
const timeRegexPattern = "\\d{2}[\\.-]\\d{2}-\\d{2}[\\.-]\\d{2}"
const timePointRegexPattern = "\\d{2}"

var client http.Client
var encoder *encoding.Encoder
var decoder *encoding.Decoder

var timeRegex *regexp.Regexp
var timePointRegex *regexp.Regexp

func init() {
	client = http.Client{
		Timeout: config.Config.RequestTimeout * time.Second,
	}

	cMap := charmap.Windows1251
	encoder = cMap.NewEncoder()
	decoder = cMap.NewDecoder()

	regex, err := regexp.Compile(timeRegexPattern)
	if err != nil {
		panic(err)
	}
	timeRegex = regex

	regex, err = regexp.Compile(timePointRegexPattern)
	if err != nil {
		panic(err)
	}
	timePointRegex = regex
}

func GetCriteria(criteriaType string) ([]models.Criterion, error) {
	request, err := http.NewRequest("GET", config.Config.TimetableUrl, nil)
	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Add("id", criteriaType)
	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	reader := decoder.Reader(response.Body)

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var result []models.Criterion
	document.Find("#vr option").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("value")
		result = append(result, models.Criterion{
			Id:   id,
			Name: s.Text(),
		})
	})

	return result, nil
}

func GetTimetable(criteriaType string, criterion string, dateRange *models.DateRange) ([]models.Day, error) {
	encodedShowArg, err := encoder.String("ПОКАЗАТЬ")
	if err != nil {
		return nil, err
	}

	formValues := url.Values{}
	formValues.Add("rel", criteriaType)
	formValues.Add("vr", criterion)
	formValues.Add("from", date_utils.ToDateString(dateRange.From))
	formValues.Add("to", date_utils.ToDateString(dateRange.To))
	formValues.Add("submit_button", encodedShowArg)

	response, err := client.PostForm(config.Config.TimetableUrl, formValues)
	if err != nil {
		return nil, err
	}

	reader := decoder.Reader(response.Body)

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	return getDays(document)
}

func getDays(document *goquery.Document) ([]models.Day, error) {
	elements := document.Find("#send td.hours")
	if elements.Length() == 0 {
		return nil, &MissingTimetableError{}
	}

	days := []models.Day{}

	// No lessons have been found
	if elements.Length() == 1 {
		return days, nil
	}

	timeRanges, err := getTimeRanges(document)
	if err != nil {
		return nil, err
	}

	i := 0
	for i < elements.Length() {
		textNode := elements.Get(i).FirstChild
		if textNode == nil || textNode.Type != html.TextNode {
			i++
			continue
		}

		date, err := date_utils.ToDate(textNode.Data)
		if err != nil {
			i++
			continue
		}

		dayOfWeek := date_utils.GetDayOfWeekName(date)

		i++

		var lessons []models.Lesson
		for ok := true; ok; ok = i < elements.Length() && !isDate(elements.Get(i)) {
			var params []string

			for j := 1; j < colCount+1; j++ {
				if elements.Get(j) != nil && elements.Get(i).FirstChild != nil {
					params = append(params, elements.Get(i).FirstChild.Data)
				} else {
					params = append(params, "")
				}
				i++
			}

			lesson := models.NewLesson(&params)

			number, err := strconv.Atoi(lesson.Number)
			if err == nil {
				lesson.Time = timeRanges[date.Weekday()][number-1]
			}

			lessons = append(lessons, *lesson)
		}

		days = append(days, models.Day{
			Date:      date_utils.ToDateString(date),
			DayOfWeek: dayOfWeek,
			Lessons:   lessons,
		})
	}

	return days, nil
}

func isDate(node *html.Node) bool {
	textNode := node.FirstChild
	if textNode == nil || textNode.Type != html.TextNode {
		return false
	}

	_, err := date_utils.ToDate(textNode.Data)
	return err == nil
}

func getTimeRanges(document *goquery.Document) (*[7][]models.TimeRange, error) {
	result := [7][]models.TimeRange{}

	rows := document.Find(".table:not(#send) tr")
	for i := 0; i < rows.Length(); i++ {
		textNodes := getAllTextNodes(rows.Get(i))
		timeRanges := *getTimeRangesFromNodes(textNodes)
		if len(timeRanges) >= 4 {
			result[0] = append(result[0], models.TimeRange{
				From: timeRanges[2].From,
				To:   timeRanges[3].To,
			})
		}
		if len(timeRanges) >= 2 {
			result[1] = append(result[1], models.TimeRange{
				From: timeRanges[0].From,
				To:   timeRanges[1].To,
			})
		}
	}

	result[6] = result[0]
	for i := 1; i < 6; i++ {
		result[i] = result[1]
	}

	return &result, nil
}

func getAllTextNodes(node *html.Node) *[]html.Node {
	var result []html.Node

	if node.Type == html.TextNode {
		result = append(result, *node)
	} else {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			childTextNodes := getAllTextNodes(child)
			result = append(result, *childTextNodes...)
		}
	}

	return &result
}

func getTimeRangesFromNodes(nodes *[]html.Node) *[]models.TimeRange {
	nodesValue := *nodes

	var result []models.TimeRange
	for i := 0; i < len(nodesValue); i++ {
		text := nodesValue[i].Data
		match := timeRegex.FindAllString(text, 1)
		if match != nil && len(match) > 0 {
			times := strings.SplitN(match[0], "-", 2)
			result = append(result, models.TimeRange{
				From: formatTime(times[0]),
				To:   formatTime(times[1]),
			})
		}
	}

	return &result
}

func formatTime(source string) string {
	timePoints := timePointRegex.FindAllString(source, 2)
	if timePoints == nil || len(timePoints) == 0 {
		return ""
	}

	result := timePoints[0]
	if len(timePoints) > 1 {
		result += ":" + timePoints[1]
	}

	return result
}

type MissingTimetableError struct{}

func (e *MissingTimetableError) Error() string {
	return "Timetable is missing on the page"
}
