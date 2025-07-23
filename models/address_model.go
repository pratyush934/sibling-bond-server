package models

type Address struct {
	Id         string `gorm:"primaryKey" json:"id"`
	UserId     string `gorm:"not null" json:"userId"`
	StreetName string `gorm:"not null" json:"streetName"`
	LandMark   string `gorm:"not null" json:"landMark"`
	ZipCode    string `gorm:"not null" json:"zipCode"`
	City       string `gorm:"not null" json:"city"`
	State      string `gorm:"not null" json:"state"`
}
