package main

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/ocassio/timetable-go-api/config"
	"github.com/ocassio/timetable-go-api/controllers"
	"log"
)

func main() {
	server := fiber.New()

	server.Use(middleware.Recover())
	server.Use(middleware.Logger())
	server.Use(cors.New())

	server.Get("/criteria/:id", controllers.GetCriteria)
	server.Get("/timetable", controllers.GetTimetable)
	server.Post("/cache/evict", controllers.EvictCache)

	err := server.Listen(config.Config.Address)
	if err != nil {
		log.Fatal(err)
	}
}
