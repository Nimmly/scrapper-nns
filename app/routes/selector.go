package routes

import (
	services "github.com/Nimmly/scrapper-nns/app/services"
	config "github.com/Nimmly/scrapper-nns/models/config"
	fiber "github.com/gofiber/fiber/v2"
)

type Response struct {
	StatusCode int      `json:"statusCode"`
	Msg        []string `json:"msg"`
}

func Checker(c *fiber.Ctx) error {
	var res Response
	selector := config.SelectorRequest{}
	if err := c.BodyParser(&selector); err != nil {
		return err
	}
	msg := services.SelectorChecker(selector)
	res = Response{
		StatusCode: 200,
		Msg:        msg,
	}
	return c.Status(res.StatusCode).JSON(res)
}
