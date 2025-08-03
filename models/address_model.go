package models

import (
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Address struct {
	Id         string `gorm:"primaryKey;type:varchar(191)" json:"id"`
	UserId     string `gorm:"not null" json:"userId"`
	StreetName string `gorm:"not null" json:"streetName"`
	LandMark   string `gorm:"not null" json:"landMark"`
	ZipCode    string `gorm:"not null" json:"zipCode"`
	City       string `gorm:"not null" json:"city"`
	State      string `gorm:"not null" json:"state"`
}

/*
Create(address *Address) (*Address, error)

GetByID(id string) (*Address, error)

GetByUserID(userID string) ([]*Address, error)

Update(address *Address) (*Address, error)

Delete(id string) error
*/

func (a *Address) BeforeCreate(t *gorm.DB) error {
	a.Id = uuid.New().String()
	return nil
}

func (a *Address) BeforeUpdate(t *gorm.DB) error {
	a.Id = uuid.New().String()
	return nil
}

func (a *Address) Create() (*Address, error) {
	if err := database.DB.Create(a).Error; err != nil {
		log.Err(err).Msg("Issue persist in Create")
		return nil, err
	}
	return a, nil
}

func GetAddressById(id string) (*Address, error) {
	var address Address
	if err := database.DB.Where(&Address{Id: id}).First(&address).Error; err != nil {
		log.Err(err).Msg("Issue persist in GetAddressById")
		return nil, err
	}
	return &address, nil
}

func GetAddressByUserId(userId string) ([]*Address, error) {
	var address []*Address
	if err := database.DB.Where(&Address{UserId: userId}).Find(&address).Error; err != nil {
		log.Err(err).Msg("Issue persist in GetAddressByUserId")
		return address, err
	}
	return address, nil
}

func UpdateAddress(address *Address) (*Address, error) {
	if err := database.DB.Save(address).Error; err != nil {
		log.Err(err).Msg("Issue persist in Update")
		return nil, err
	}
	return address, nil
}

func DeleteAddress(id string) error {
	return database.DB.Where(&Address{Id: id}).Delete(&Address{}).Error
}
