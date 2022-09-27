package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Model
	Email      string       `gorm:"uniqueindex;not null" json:"email"`
	Password   string       `gorm:"not null"`
	Role       string       `gorm:"type:varchar(255);not null" json:"role"`
	VerifiedAt sql.NullTime `gorm:"index"`
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave(db *gorm.DB) error {
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) InsertUser(db *gorm.DB) (*User, error) {
	err := db.Create(&u).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) GetUserByID(db *gorm.DB, id uint) (*User, error) {
	err := db.First(&u, id).Error

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	err := db.Where("email = ?").First(&u).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}
