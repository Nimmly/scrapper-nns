package main

import (
	app "github.com/Nimmly/scrapper-nns/app"
	fiber "github.com/gofiber/fiber/v2"
)

func main() {
	a := fiber.New()
	app.SetupRoutes(a)

	a.Listen(":1337")
}
