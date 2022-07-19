package models

type City struct {
	ID          uint   `gorm:"primaryKey"`
	AddressBody string `json:"addressBody"`
	Counts      int    `json:"counts"`
}
