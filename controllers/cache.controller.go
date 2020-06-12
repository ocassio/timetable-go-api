package controllers

import (
	"github.com/gofiber/fiber"
	"github.com/ocassio/timetable-go-api/services/data_provider"
)

func EvictCache(ctx *fiber.Ctx) {
	data_provider.EvictCache()
	ctx.SendStatus(200)
}
