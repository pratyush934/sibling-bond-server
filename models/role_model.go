package models

type Role struct {
	Id          int    `gorm:"primaryKey" json:"id"`
	RoleName    string `gorm:"not null" json:"roleName"`
	Description string `gorm:"not null" json:"description"`
}

/*
	Role Id 1 is for User,
	Role Id 2 is for admin
*/
