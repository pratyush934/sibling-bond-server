package models

import (
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Image struct {
	Id        string `gorm:"primaryKey;type:varchar(191)" json:"id"`
	URL       string `gorm:"not null" json:"url"`
	FileName  string `gorm:"not null" json:"fileName"`
	FieldId   string `gorm:"not null" json:"fieldId"`
	AltText   string `json:"altText"`
	IsPrimary bool   `gorm:"default:false" json:"isPrimary"`
	SortOrder int    `gorm:"default:0" json:"sortOrder"`
	FileSize  int64  `json:"fileSize"`
	MimeType  string `json:"mimeType"`

	// Foreign key
	ProductId string `gorm:"not null;type:varchar(191)" json:"productId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (i *Image) BeforeCreate(tx *gorm.DB) error {
	i.Id = uuid.New().String()
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Image) BeforeUpdate(tx *gorm.DB) error {
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Image) CreateImage() (*Image, error) {
	if err := database.DB.Create(i).Error; err != nil {
		log.Err(err).Msg("Issue exist while CreatingImage")
		return nil, err
	}
	return i, nil
}
