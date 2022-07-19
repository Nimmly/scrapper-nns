package models

type Image struct {
	ID          uint `gorm:"primaryKey"`
	ApartmentID string
	Link        string
}
type ImageResponse struct {
	ApartmentID string
	Link        string
}

// func CreateResponseImages(image Image) ImageResponse {
// 	return Image{ApartmentID: image.ApartmentID, Link: image.Link}
// }
