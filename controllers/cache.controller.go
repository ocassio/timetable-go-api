package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ocassio/timetable-go-api/services/data_provider"
)

func EvictCache(ctx *fiber.Ctx) error {
	data_provider.EvictCache()
	ctx.SendStatus(200)
	return nil
}
