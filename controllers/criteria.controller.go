package controllers

import (
	"github.com/gofiber/fiber"
	"github.com/ocassio/timetable-go-api/services/data_provider"
)

func GetCriteria(ctx *fiber.Ctx) {
	criteriaType := ctx.Params("id")
	criteria, err := data_provider.GetCriteria(criteriaType)
	if err != nil {
		ctx.Next(err)
	}

	err = ctx.JSON(criteria)
	if err != nil {
		ctx.Next(err)
	}
}
