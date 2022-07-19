package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	database "github.com/Nimmly/scrapper-nns/database"
	models "github.com/Nimmly/scrapper-nns/models"
	fiber "github.com/gofiber/fiber/v2"
)

type City struct {
	AddressBody string `json:"addressBody"`
	Counts      int    `json:"counts"`
}

func CreateResponseCity(city models.City) City {
	return City{AddressBody: city.AddressBody, Counts: city.Counts}
}

func ListCitiesNekretnine(c *fiber.Ctx) error {
	var city []City
	var ct City

	response, _ := http.Get("https://www.nekretnine.rs/app-api/get-cities?parent=nqzAYmYBN5vgk7Geu8TB")
	body, _ := ioutil.ReadAll(response.Body)
	err1 := json.Unmarshal(body, &city)
	if err1 != nil {
		log.Fatal(err1)
	}
	for _, c := range city {
		if c.Counts > 0 {
			short := strings.SplitAfter(c.AddressBody, ",")
			urlReady := strings.Replace(short[0], ",", "", 1)
			ready := strings.Replace(urlReady, " ", "-", 1)
			readyLower := strings.ToLower(ready)
			unicode := []string{"ž", "ć", "č", "đ", "š"}
			sub := []string{"z", "c", "c", "dj", "s"}

			var temp string

			for idx, letter := range unicode {
				temp = strings.ReplaceAll(readyLower, letter, sub[idx])
				if temp != readyLower {
					readyLower = temp
				}
				ct.AddressBody = readyLower
			}

			ct.Counts = c.Counts
			database.Database.Db.Create(&ct)
		}
	}
	return c.Status(200).JSON(ct)
}
