package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	services "github.com/Nimmly/scrapper-nns/app/services"
	database "github.com/Nimmly/scrapper-nns/database"
	models "github.com/Nimmly/scrapper-nns/models"
	config "github.com/Nimmly/scrapper-nns/models/config"
	fiber "github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func GetAllApartments(c *fiber.Ctx) error {
	var apartments []models.ApartmentResponse
	var responseApapartments []models.ApartmentResponse
	limit := c.Query("limit")
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Fatal(err)
	}
	err1 := database.Database.Db.Model(&models.Apartment{}).Limit(intLimit).Preload("Images").Find(&apartments).Error
	if err1 != nil {
		fmt.Println(err1)
	}
	for _, apartment := range apartments {
		singleFlat := models.CreateResponseApartment(apartment)
		responseApapartments = append(responseApapartments, singleFlat)
	}
	return c.Status(200).JSON(responseApapartments)
}

func GeneralScrapper(c *fiber.Ctx) error {
	var response config.Success
	NekretnineRS := config.SiteMapper{
		MainSite:               "https://www.nekretnine.rs",
		ScrapeLink:             "https://www.nekretnine.rs/stambeni-objekti/stanovi/izdavanje-prodaja/izdavanje/lista/po-stranici/20/",
		NextPage:               ".next-number",
		ApartmentListContainer: ".advert-list",
		ContainerSingleItem:    ".row",
		SingleItemLink:         "h2 > a",
		PreviewContainer:       ".property__body",
		PreviewTitle:           ".detail-title",
		PreviewDescription:     ".cms-content-inner",
		PreviewLocation:        ".stickyBox__Location",
		PreviewSize:            ".stickyBox__size",
		PreviewPrice:           ".stickyBox__price",
		Form:                   true,
		FormURL:                "https://www.nekretnine.rs/form/show-phone-number/phone/",
		FormMethod:             "POST",
		FormHeaders:            []string{"X-Requested-With", "XMLHttpRequest"},
		FormAttr:               `[data-dynamic-form="show-phone-number"]`,
		FormType:               "action",
		FormDelimiter:          "/",
		DelimiterPositions:     []int{0, 4},
		Category:               "rent",
	}
	services.Scrapper(NekretnineRS)
	response = config.Success{StatusCode: 201, Msg: "Testing apartment scrape Completed!"}
	return c.Status(200).JSON(response)
}
func ScrapperWithRequest(c *fiber.Ctx) error {
	var req config.SiteMapper
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	fmt.Println(&req)
	services.Scrapper(req)
	return c.Status(200).JSON(map[string]string{"msg": "Testing apartment scrape Completed!"})
}
func ImageScrapper(c *fiber.Ctx) error {
	var response config.Success
	NekretnineRS := config.ImagesSiteMapper{
		ImagesURLList:           ".advert-list",
		SingleImageRow:          ".row",
		ImageURLAttr:            "h2 > a",
		MainSite:                "https://www.nekretnine.rs",
		Fragment:                "galerija",
		Carousel:                ".carousel",
		FlatIdDelimiter:         "/",
		FlatIdDelimiterPosition: 6,
		ScrapeLink:              "https://www.nekretnine.rs/stambeni-objekti/stanovi/izdavanje-prodaja/izdavanje/lista/po-stranici/20/",
	}
	services.ScrappeImages(NekretnineRS)
	response = config.Success{StatusCode: 201, Msg: "Testing Image Scrapper Completed!"}
	return c.Status(200).JSON(response)
}
func Rabbit(c *fiber.Ctx) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")
	data := config.SiteMapper{
		MainSite:               "https://www.nekretnine.rs",
		ScrapeLink:             "https://www.nekretnine.rs/stambeni-objekti/stanovi/izdavanje-prodaja/izdavanje/lista/po-stranici/20/",
		NextPage:               ".next-number",
		ApartmentListContainer: ".advert-list",
		ContainerSingleItem:    ".row",
		SingleItemLink:         "h2 > a",
		PreviewContainer:       ".property__body",
		PreviewTitle:           ".detail-title",
		PreviewDescription:     ".cms-content-inner",
		PreviewLocation:        ".stickyBox__Location",
		PreviewSize:            ".stickyBox__size",
		PreviewPrice:           ".stickyBox__price",
		Form:                   true,
		FormURL:                "https://www.nekretnine.rs/form/show-phone-number/phone/",
		FormMethod:             "POST",
		FormHeaders:            []string{"X-Requested-With", "XMLHttpRequest"},
		FormAttr:               `[data-dynamic-form="show-phone-number"]`,
		FormType:               "action",
		FormDelimiter:          "/",
		DelimiterPositions:     []int{0, 4},
		Category:               "rent",
	}
	t, _ := json.Marshal(data)
	fmt.Println(database.Database)
	body := t
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
	return c.SendString("Scrape in Progress!")
}
