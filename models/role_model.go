package models

import "gorm.io/gorm"

type Role struct {
	Id          int    `gorm:"primaryKey" json:"id"`
	RoleName    string `gorm:"not null" json:"roleName"`
	Description string `gorm:"not null" json:"description"`
}

/*
	Role Id 1 is for User,
	Role Id 2 is for admin
	Role Id 3 for tenant
*/

func CreateRole(db *gorm.DB, role *Role) error {
	return db.Create(role).Error
}

func GetRoleByID(db *gorm.DB, id int) (*Role, error) {
	var role Role
	err := db.First(&role, id).Error
	return &role, err
}

func GetAllRoles(db *gorm.DB) ([]Role, error) {
	var roles []Role
	err := db.Find(&roles).Error
	return roles, err
}

func UpdateRole(db *gorm.DB, role *Role) error {
	return db.Save(role).Error
}

func DeleteRole(db *gorm.DB, id int) error {
	return db.Delete(&Role{}, id).Error
}
