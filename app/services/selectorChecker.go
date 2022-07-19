package services

import (
	"fmt"

	config "github.com/Nimmly/scrapper-nns/models/config"
	colly "github.com/gocolly/colly"
)

type SelectorRequest struct {
	NextPage            bool
	NextPageSelector    string
	ListOfItems         bool
	ListOfItemsSelector string
	SingleItemSelector  string
	SingleItemLink      string
	MainSite            string
	ScrapeLink          string
	SingleItemPage      bool
	SinglePageContainer string
	SinglePageSelector  string
}

func SelectorChecker(params config.SelectorRequest) []string {
	var nextPageSelector string
	var singleItemLink string
	var singlePageSelector string

	cr := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)
	cr.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})
	cr.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})
	if params.NextPage {
		cr.OnHTML(params.NextPageSelector, func(e *colly.HTMLElement) {
			nextPageSelector = e.Request.AbsoluteURL(e.Attr("href"))
		})
	}
	if params.ListOfItems {
		cr.OnHTML(params.ListOfItemsSelector, func(e *colly.HTMLElement) {
			e.ForEach(params.SingleItemSelector, func(_ int, el *colly.HTMLElement) {
				link := el.ChildAttr(params.SingleItemLink, "href")
				singleItemLink = params.MainSite + link
			})
		})
	}
	if params.SingleItemPage {
		cr.OnHTML(params.SinglePageContainer, func(e *colly.HTMLElement) {
			// var links []string
			singlePageSelector = e.ChildText(params.SinglePageSelector)

			// form := e.ChildAttr(site.FormAttr, site.FormType)
			// links = append(links, form)
			// delimiter := site.FormDelimiter
			// siteId := strings.Split(links[site.DelimiterPositions[0]], delimiter)[site.DelimiterPositions[1]]
			// apartment.SiteID = siteId
		})
	}
	cr.Visit(params.ScrapeLink)
	cr.Wait()

	//
	// apartments := []models.Apartment{}
	// database.Database.Db.Find(&apartments)
	// var phone models.PhoneResponse
	// for _, flat := range apartments {
	// 	link := site.FormURL + flat.SiteID
	// 	req, err := http.NewRequest(site.FormMethod, link, http.NoBody)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	req.Header.Set(site.FormHeaders[0], site.FormHeaders[1])
	// 	client := &http.Client{}
	// 	res, e := client.Do(req)
	// 	if e != nil {
	// 		log.Fatal(e)
	// 	}
	// 	body, _ := ioutil.ReadAll(res.Body)
	// 	err1 := json.Unmarshal(body, &phone)
	// 	if err1 != nil {
	// 		log.Fatal(err1)
	// 	}
	// 	database.Database.Db.Table("apartments").Where("site_id = ?", flat.SiteID).Update("phone", phone.Number[0])
	// }
	// //
	// res := map[string]string{}
	return []string{nextPageSelector, singleItemLink, singlePageSelector}
}
