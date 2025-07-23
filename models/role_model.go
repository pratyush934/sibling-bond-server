package models

type Role struct {
	Id          int    `gorm:"primaryKey" json:"id"`
	RoleName    string `gorm:"not null" json:"roleName"`
	Description string `gorm:"not null" json:"description"`
}
