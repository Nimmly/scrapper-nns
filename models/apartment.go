package models

type Apartment struct {
	ID       uint   `gorm:"primaryKey"`
	Title    string `gorm:"index"`
	Location string `json:"location"`
	Price    string `json:"price"`
	Size     string `json:"size"`
	Desc     string `json:"desc"`
	Phone    string `json:"phone"`
	SiteID   string `gorm:"unique"`
	Category string `json:"category"`
}

type ApartmentResponse struct {
	Title    string  `json:"title"`
	Location string  `json:"location"`
	Price    string  `json:"price"`
	Size     string  `json:"size"`
	Desc     string  `json:"desc"`
	Phone    string  `json:"phone"`
	SiteID   string  `json:"siteId"`
	Category string  `json:"category"`
	Images   []Image `json:"images" gorm:"foreignKey:ApartmentID;references:SiteID"`
}

type PhoneResponse struct {
	Success             bool     `json:"success"`
	Number              []string `json:"phone"`
	AdvertiserGaEventId string   `json:"advertiserGaEventId"`
}

func CreateResponseApartment(apartment ApartmentResponse) ApartmentResponse {
	return ApartmentResponse{Title: apartment.Title, Location: apartment.Location, Price: apartment.Price, Size: apartment.Size, Desc: apartment.Desc, Phone: apartment.Phone, SiteID: apartment.SiteID, Category: apartment.Category, Images: apartment.Images}
}
