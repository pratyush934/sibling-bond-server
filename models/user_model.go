package models

import (
	"github.com/google/uuid"
	"github.com/pratyush934/sibling-bond-server/database"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id                  string     `json:"id" gorm:"primaryKey; type:varchar(100)"`
	Email               string     `gorm:"unique;not null" json:"email"`
	PassWord            string     `gorm:"not null" json:"passWord"`
	UserName            string     `gorm:"not null;unique" json:"userName"`
	FirstName           string     `gorm:"not null" json:"firstName"`
	LastName            string     `json:"lastName"`
	PhoneNumber         string     `json:"phoneNumber"`
	RoleId              int        `gorm:"not null; default:1" json:"roleId"`
	Role                Role       `gorm:"constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"role"`
	Addresses           []Address  `gorm:"foreignKey:UserId" json:"addresses"`
	Orders              []Order    `gorm:"foreignKey:UserId" json:"orders"`
	PrimaryAddress      string     `json:"primaryAddress"`
	PasswordResetToken  *string    `json:"-"`
	PasswordResetExpiry *time.Time `json:"-"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

func (u *User) BeforeCreate(t *gorm.DB) error {
	u.Id = uuid.New().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	password, err := bcrypt.GenerateFromPassword([]byte(u.PassWord), bcrypt.DefaultCost)
	if err != nil {
		log.Err(err).Msg("Issue exist in BeforeCreate here part 1")
		return err
	}
	u.PassWord = string(password)
	lastName := u.LastName
	if lastName == "" {
		lastName = "user"
	}
	u.UserName = u.FirstName + "." + lastName + "." + uuid.New().String()

	return nil
}

func (u *User) BeforeUpdate(t *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ValidatePassWord(pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PassWord), []byte(pass))
	return err == nil
}

func (u *User) CreateUser() (*User, error) {
	if err := database.DB.Create(u).Error; err != nil {
		log.Err(err).Msg("Issue while creating User")
		return nil, err
	}
	return u, nil
}

func (u *User) GeneratePasswordResetToken() string {
	token := uuid.New().String()
	expiryTime := time.Now().Add(1 * time.Hour)

	u.PasswordResetToken = &token
	u.PasswordResetExpiry = &expiryTime

	return token
}

func (u *User) ValidatePasswordResetToken(token string) bool {
	if u.PasswordResetToken == nil || u.PasswordResetExpiry == nil {
		return false
	}
	return *u.PasswordResetToken == token && time.Now().Before(*u.PasswordResetExpiry)
}

func (u *User) ResetPassword(newPassword string) error {
	password, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassWord = string(password)
	u.PasswordResetExpiry = nil
	u.PasswordResetToken = nil

	return nil
}

func GetUserByResetToken(token string) (*User, error) {
	var user User
	if err := database.DB.Where("password_reset_token = ? AND password_reset_expiry > ?", token, time.Now()).First(&user).Error; err != nil {
		log.Err(err).Msg("Invalid or expired reset token")
		return nil, err
	}
	return &user, nil
}

func InitiatePasswordReset(email string) (*User, string, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}

	token := user.GeneratePasswordResetToken()

	if err := database.DB.Save(user).Error; err != nil {
		log.Err(err).Msg("Failed to save password reset token")
		return nil, "", err
	}

	return user, token, nil
}

func GetUserById(id string) (*User, error) {
	var user User
	if err := database.DB.Where(&User{Id: id}).First(&user).Error; err != nil {
		log.Err(err).Msg("Didn't get the User")
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := database.DB.Where(&User{Email: email}).First(&user).Error; err != nil {
		log.Err(err).Msg("Didn't get the user by email")
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *User) (*User, error) {
	if err := database.DB.Updates(user).Error; err != nil {
		log.Err(err).Msg("issue while updating the user")
		return &User{}, err
	}
	return user, nil
}

func DeleteUser(id string) error {
	if err := database.DB.Where(&User{Id: id}).Delete(&User{}).Error; err != nil {
		log.Err(err).Msg("Issue while deleting the user")
		return err
	}
	return nil
}

func GetAllUsers(offset, limit int) ([]User, error) {
	var users []User
	if err := database.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		log.Err(err).Msg("Issue while getting all the users")
		return users, err
	}
	return users, nil
}
