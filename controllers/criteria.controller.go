package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ocassio/timetable-go-api/services/data_provider"
)

func GetCriteria(ctx *fiber.Ctx) error {
	criteriaType := ctx.Params("id")
	criteria, err := data_provider.GetCriteria(criteriaType)
	if err != nil {
		return err
	}

	err = ctx.JSON(criteria)
	if err != nil {
		return err
	}

	return nil
}
