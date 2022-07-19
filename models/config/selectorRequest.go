package config

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
