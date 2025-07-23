package models

import "time"

type User struct {
	Id             string    `json:"id" gorm:"primaryKey; type:varchar(100)"`
	Email          string    `gorm:"unique; not null" json:"email"`
	PassWord       string    `gorm:"not null" json:"passWord"`
	FirstName      string    `gorm:"not null" json:"firstName"`
	LastName       string    `json:"lastName"`
	PhoneNumber    string    `json:"phoneNumber"`
	RoleId         int       `gorm:"not null default:1" json:"roleId"`
	Role           Role      `gorm:"constraint:onUpdate:CASCADE onDelete:CASCADE" json:"role"`
	Addresses      []Address `gorm:"foreignKey:UserId" json:"addresses"`
	Orders         []Order   `gorm:"foreignKey:UserId" json:"orders"`
	PrimaryAddress string    `json:"primaryAddress"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
