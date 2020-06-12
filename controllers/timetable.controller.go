package controllers

import (
	"github.com/gofiber/fiber"
	"github.com/ocassio/timetable-go-api/models"
	"github.com/ocassio/timetable-go-api/services/data_provider"
	"github.com/ocassio/timetable-go-api/utils/date_utils"
)

func GetTimetable(ctx *fiber.Ctx) {
	criteriaType := ctx.Query("criteriaType")
	criterion := ctx.Query("criterion")
	from := ctx.Query("from")
	to := ctx.Query("to")

	var dateRange models.DateRange
	if len(from) == 0 {
		dateRange = date_utils.GetSevenDays(nil)
	} else {
		fromDate, err := date_utils.ToDate(from)
		if err != nil {
			sendMalformedDateError(ctx, from)
			return
		}

		if len(to) > 0 {
			toDate, err := date_utils.ToDate(to)
			if err != nil {
				sendMalformedDateError(ctx, to)
				return
			}

			dateRange = models.DateRange{
				From: fromDate,
				To:   toDate,
			}
		} else {
			dateRange = date_utils.GetSevenDays(&fromDate)
		}
	}

	timetable, err := data_provider.GetLessons(criteriaType, criterion, &dateRange)
	if err != nil {
		ctx.Next(err)
	}

	err = ctx.JSON(timetable)
	if err != nil {
		ctx.Next(err)
	}
}

func sendMalformedDateError(ctx *fiber.Ctx, date string) {
	_ = ctx.Status(400).JSON(fiber.Map{
		"error": "Malformed date: " + date,
	})
}
