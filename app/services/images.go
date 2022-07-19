package services

import (
	"fmt"
	"strings"

	database "github.com/Nimmly/scrapper-nns/database"
	models "github.com/Nimmly/scrapper-nns/models"
	config "github.com/Nimmly/scrapper-nns/models/config"
	colly "github.com/gocolly/colly"
)

func ScrappeImages(site config.ImagesSiteMapper) string {
	var image models.ImageResponse
	cr := colly.NewCollector(
		colly.MaxDepth(3),
		colly.Async(true),
	)
	cr.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})
	cr.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})
	// cr.OnHTML(".next-number", func(e *colly.HTMLElement) {
	// 	nextPage := e.Request.AbsoluteURL(e.Attr("href"))
	// 	cr.Visit(nextPage)
	// })

	cr.OnHTML(site.ImagesURLList, func(e *colly.HTMLElement) {
		e.ForEach(site.SingleImageRow, func(_ int, el *colly.HTMLElement) {
			link := el.ChildAttr(site.ImageURLAttr, "href")
			cr.Visit(site.MainSite + link + site.Fragment)
		})
	})
	cr.OnHTML(site.Carousel, func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		del := site.FlatIdDelimiter
		siteId := strings.Split(url, del)[site.FlatIdDelimiterPosition]
		e.ForEach("img", func(_ int, el *colly.HTMLElement) {
			link := el.Attr("src")
			if link != "" {
				image.ApartmentID = siteId
				image.Link = link
				database.Database.Db.Table("images").Create(&image)
			}
		})
	})
	cr.Visit(site.ScrapeLink)
	cr.Wait()
	return "success"
}
