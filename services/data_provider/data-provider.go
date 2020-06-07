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
	"time"
)

const colCount = 7

var client http.Client
var encoder *encoding.Encoder
var decoder *encoding.Decoder

func init() {
	client = http.Client {
		Timeout: config.Config.RequestTimeout * time.Second,
	}

	cMap := charmap.Windows1251
	encoder = cMap.NewEncoder()
	decoder = cMap.NewDecoder()
}


func GetCriteria(criteriaType string) ([]models.Criterion, error) {
	request, err := http.NewRequest("GET", config.Config.TimetableUrl, nil); if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Add("id", criteriaType)
	request.URL.RawQuery = query.Encode()

	response, err := client.Do(request); if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	reader := decoder.Reader(response.Body)

	document, err := goquery.NewDocumentFromReader(reader); if err != nil {
		return nil, err
	}

	var result []models.Criterion
	document.Find("#vr option").Each(func (i int, s *goquery.Selection) {
		id, _ := s.Attr("value")
		result = append(result, models.Criterion {
			Id: id,
			Name: s.Text(),
		})
	})

	return result, nil
}

func GetTimetable(criteriaType string, criterion string, dateRange *models.DateRange) ([]models.Day, error) {
	encodedShowArg, err := encoder.String("ПОКАЗАТЬ"); if err != nil {
		return nil, err
	}

	formValues := url.Values {}
	formValues.Add("rel", criteriaType)
	formValues.Add("vr", criterion)
	formValues.Add("from", date_utils.ToDateString(dateRange.From))
	formValues.Add("to", date_utils.ToDateString(dateRange.To))
	formValues.Add("submit_button", encodedShowArg)

	response, err := client.PostForm(config.Config.TimetableUrl, formValues); if err != nil {
		return nil, err
	}

	reader := decoder.Reader(response.Body)

	document, err := goquery.NewDocumentFromReader(reader); if err != nil {
		return nil, err
	}

	elements := document.Find("#send td.hours")
	if elements.Length() == 0 {
		return nil, &MissingTimetableError{}
	}

	return getDays(elements), nil
}

func getDays(elements *goquery.Selection) []models.Day {
	days := []models.Day{}

	// No lessons have been found
	if elements.Length() == 1 {
		return days
	}

	i := 0
	for ; i < elements.Length(); {
		textNode := elements.Get(i).FirstChild
		if textNode == nil || textNode.Type != html.TextNode {
			i++; continue
		}

		date, err := date_utils.ToDate(textNode.Data); if err != nil {
			i++; continue
		}

		dayOfWeek := date_utils.GetDayOfWeekName(date)

		i++

		var lessons []models.Lesson
		for ok := true; ok; ok = i < elements.Length() && !isDate(elements.Get(i)) {
			var params []string

			for j := 0; j < colCount; j++ {
				params = append(params, elements.Get(j).FirstChild.Data)
				i++
			}

			lesson := models.NewLesson(&params)
			lessons = append(lessons, *lesson)
		}

		days = append(days, models.Day {
			Date: date_utils.ToDateString(date),
			DayOfWeek: dayOfWeek,
			Lessons: lessons,
		})
	}

	return days
}

func isDate(node *html.Node) bool {
	textNode := node.FirstChild
	if textNode == nil || textNode.Type != html.TextNode {
		return false
	}

	_, err := date_utils.ToDate(textNode.Data)
	return err == nil
}

type MissingTimetableError struct {}

func (e *MissingTimetableError) Error() string {
	return "Timetable is missing on the page"
}
