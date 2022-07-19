package app

import (
	routes "github.com/Nimmly/scrapper-nns/app/routes"
	database "github.com/Nimmly/scrapper-nns/database"
	fiber "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	database.ConnectDb()
	api := app.Group("/api")
	// General
	api.Get("/get-all", routes.GetAllApartments)

	// Scrapper
	api.Get("/general_scrapper", routes.GeneralScrapper)
	api.Get("/general_image_scrapper", routes.ImageScrapper)

	// Checkers
	api.Post("/check", routes.Checker)

	// Testing message broker RabbitMQ
	api.Get("/rabbit", routes.Rabbit)

	// Testing with Request.Body
	api.Post("/dynamic_scrapper", routes.ScrapperWithRequest)
}
