package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	models "github.com/Nimmly/scrapper-nns/models"
	config "github.com/Nimmly/scrapper-nns/models/config"
	colly "github.com/gocolly/colly"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type DB struct {
	Db *gorm.DB
}

var Database DB

func ConnectDatabase() {
	dbinfo := "user=nimmly password=root host=localhost port=5432 dbname=nns sslmode=disable"
	db, err := gorm.Open(postgres.Open(dbinfo))
	if err != nil {
		log.Fatal("Postgres instance terminated... Check your connection!")
	}
	log.Println("Connected to the database")
	db.Logger = logger.Default.LogMode(logger.Info)
	Database = DB{Db: db}
}

func main() {
	ConnectDatabase()
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

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var selectors config.SiteMapper
			json.Unmarshal(d.Body, &selectors)

			var apartment models.ApartmentResponse

			cr := colly.NewCollector(
				colly.MaxDepth(2),
				colly.Async(true),
			)
			cr.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})
			cr.OnRequest(func(r *colly.Request) {
				fmt.Printf("Visiting %s\n", r.URL)
			})
			cr.OnHTML(selectors.NextPage, func(e *colly.HTMLElement) {
				nextPage := e.Request.AbsoluteURL(e.Attr("href"))
				cr.Visit(nextPage)
			})

			cr.OnHTML(selectors.ApartmentListContainer, func(e *colly.HTMLElement) {
				e.ForEach(selectors.ContainerSingleItem, func(_ int, el *colly.HTMLElement) {
					link := el.ChildAttr(selectors.SingleItemLink, "href")
					cr.Visit(selectors.MainSite + link)
				})
			})
			cr.OnHTML(selectors.PreviewContainer, func(e *colly.HTMLElement) {
				var links []string
				apartment.Title = strings.TrimSpace(e.ChildText(selectors.PreviewTitle))
				apartment.Location = strings.TrimSpace(e.ChildText(selectors.PreviewLocation))
				apartment.Price = strings.TrimSpace(e.ChildText(selectors.PreviewPrice))
				apartment.Size = strings.TrimSpace(e.ChildText(selectors.PreviewSize))
				apartment.Desc = strings.TrimSpace(e.ChildText(selectors.PreviewDescription))

				if selectors.Form {
					form := e.ChildAttr(selectors.FormAttr, selectors.FormType)
					links = append(links, form)
					delimiter := selectors.FormDelimiter
					siteId := strings.Split(links[selectors.DelimiterPositions[0]], delimiter)[selectors.DelimiterPositions[1]]
					apartment.SiteID = siteId
				}
				apartment.Category = selectors.Category
				Database.Db.Table("apartments").Create(&apartment)

			})
			cr.Visit(selectors.ScrapeLink)
			cr.Wait()

			//
			apartments := []models.Apartment{}
			Database.Db.Find(&apartments)
			var phone models.PhoneResponse
			for _, flat := range apartments {
				link := selectors.FormURL + flat.SiteID
				req, err := http.NewRequest(selectors.FormMethod, link, http.NoBody)
				if err != nil {
					log.Fatal(err)
				}
				req.Header.Set(selectors.FormHeaders[0], selectors.FormHeaders[1])
				client := &http.Client{}
				res, e := client.Do(req)
				if e != nil {
					log.Fatal(e)
				}
				body, _ := ioutil.ReadAll(res.Body)
				err1 := json.Unmarshal(body, &phone)
				if err1 != nil {
					log.Fatal(err1)
				}
				Database.Db.Table("apartments").Where("site_id = ?", flat.SiteID).Update("phone", phone.Number[0])
			}
			log.Printf("Received a message: %s", d.Body)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
