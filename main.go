package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ocassio/timetable-go-api/config"
	"github.com/ocassio/timetable-go-api/controllers"
)

func main() {
	server := fiber.New()

	server.Use(recover.New())
	server.Use(logger.New())
	server.Use(cors.New())

	server.Get("/criteria/:id", controllers.GetCriteria)
	server.Get("/timetable", controllers.GetTimetable)
	server.Post("/cache/evict", controllers.EvictCache)

	err := server.Listen(config.Config.Address)
	if err != nil {
		log.Fatal(err)
	}
}
