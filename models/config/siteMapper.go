package config

type SiteMapper struct {
	MainSite               string
	ScrapeLink             string
	NextPage               string
	ApartmentListContainer string
	ContainerSingleItem    string
	SingleItemLink         string
	PreviewContainer       string
	PreviewTitle           string
	PreviewDescription     string
	PreviewLocation        string
	PreviewSize            string
	PreviewPrice           string
	Form                   bool
	FormURL                string
	FormMethod             string
	FormHeaders            []string
	FormAttr               string
	FormType               string
	FormDelimiter          string
	DelimiterPositions     []int
	Category               string
}
