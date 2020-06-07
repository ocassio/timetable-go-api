package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ocassio/timetable-go-api/config"
	"github.com/ocassio/timetable-go-api/models"
	"github.com/ocassio/timetable-go-api/services/data_provider"
	"github.com/ocassio/timetable-go-api/utils/date_utils"
	"log"
	"net/http"
)

func main() {
	server := gin.Default()

	server.GET("/criteria/:id", func (ctx *gin.Context) {
		criteriaType := ctx.Param("id")
		criteria, err := data_provider.GetCriteria(criteriaType); if err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusOK, criteria)
	})

	server.GET("/timetable", func(ctx *gin.Context) {
		criteriaType := ctx.Query("criteriaType")
		criterion := ctx.Query("criterion")
		from, fromPresent := ctx.GetQuery("from")
		to, toPresent := ctx.GetQuery("to")

		var dateRange models.DateRange
		if !fromPresent {
			dateRange = date_utils.GetSevenDays(nil)
		} else {
			fromDate, err := date_utils.ToDate(from); if err != nil {
				sendMalformedDateError(ctx, from)
				return
			}

			if toPresent {
				toDate, err := date_utils.ToDate(to); if err != nil {
					sendMalformedDateError(ctx, to)
					return
				}

				dateRange = models.DateRange {
					From: fromDate,
					To: toDate,
				}
			} else {
				dateRange = date_utils.GetSevenDays(&fromDate)
			}
		}

		timetable, err := data_provider.GetTimetable(criteriaType, criterion, &dateRange); if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, timetable)
	})

	err := server.Run(config.Config.Address); if err != nil {
		log.Fatal(err)
	}
}

func sendMalformedDateError(ctx *gin.Context, date string) {
	ctx.AbortWithStatusJSON(400, gin.H {
		"error": "Malformed date: " + date,
	})
}
