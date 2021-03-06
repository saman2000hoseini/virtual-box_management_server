package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
}

func NewUser(username, password, role string) *User {
	return &User{
		Username: username,
		Password: password,
		Role:     role,
	}
}

func (u *User) CheckPassword(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	return err == nil
}

type UserRepo interface {
	Find(username string) (User, error)
	Save(user *User) error
	Update(user User) error
}

type SQLUserRepo struct {
	DB *gorm.DB
}

func (r SQLUserRepo) Find(username string) (User, error) {
	var stored User
	err := r.DB.Where(&User{Username: username}).First(&stored).Error

	return stored, err
}

func (r SQLUserRepo) Save(user *User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	return r.DB.Create(user).Error
}

func (r SQLUserRepo) Update(user User) error {
	return r.DB.Save(user).Error
}
