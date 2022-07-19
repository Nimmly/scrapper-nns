package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	database "github.com/Nimmly/scrapper-nns/database"
	models "github.com/Nimmly/scrapper-nns/models"
	config "github.com/Nimmly/scrapper-nns/models/config"
	colly "github.com/gocolly/colly"
)

func Scrapper(site config.SiteMapper) string {
	var apartment models.ApartmentResponse

	cr := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)
	cr.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})
	cr.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})
	cr.OnHTML(site.NextPage, func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		cr.Visit(nextPage)
	})

	cr.OnHTML(site.ApartmentListContainer, func(e *colly.HTMLElement) {
		e.ForEach(site.ContainerSingleItem, func(_ int, el *colly.HTMLElement) {
			link := el.ChildAttr(site.SingleItemLink, "href")
			cr.Visit(site.MainSite + link)
		})
	})
	cr.OnHTML(site.PreviewContainer, func(e *colly.HTMLElement) {
		var links []string
		apartment.Title = strings.TrimSpace(e.ChildText(site.PreviewTitle))
		apartment.Location = strings.TrimSpace(e.ChildText(site.PreviewLocation))
		apartment.Price = strings.TrimSpace(e.ChildText(site.PreviewPrice))
		apartment.Size = strings.TrimSpace(e.ChildText(site.PreviewSize))
		apartment.Desc = strings.TrimSpace(e.ChildText(site.PreviewDescription))

		//
		if site.Form {
			form := e.ChildAttr(site.FormAttr, site.FormType)
			links = append(links, form)
			delimiter := site.FormDelimiter
			siteId := strings.Split(links[site.DelimiterPositions[0]], delimiter)[site.DelimiterPositions[1]]
			apartment.SiteID = siteId
		} else {
			url := e.Request.URL.String()
			del := site.FormDelimiter
			siteId := strings.Split(url, del)[site.DelimiterPositions[0]]
			apartment.SiteID = siteId
		}
		apartment.Category = site.Category
		//
		database.Database.Db.Table("apartments").Create(&apartment)
	})
	cr.Visit(site.ScrapeLink)
	cr.Wait()

	//
	apartments := []models.Apartment{}
	database.Database.Db.Find(&apartments)
	var phone models.PhoneResponse
	for _, flat := range apartments {
		link := site.FormURL + flat.SiteID
		req, err := http.NewRequest(site.FormMethod, link, http.NoBody)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set(site.FormHeaders[0], site.FormHeaders[1])
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
		database.Database.Db.Table("apartments").Where("site_id = ?", flat.SiteID).Update("phone", phone.Number[0])
	}
	//
	return "success"
}
