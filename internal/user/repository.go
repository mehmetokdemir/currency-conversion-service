package user

import (
	// Go imports
	"errors"
	"fmt"

	// External imports
	_ "github.com/golang/mock/mockgen/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user User) (*User, error)
	IsUserExistWithSameUsername(username string) bool
	IsUserExistWithSameEmail(email string) bool
	GetUserByUsername(username string) (*User, error)
	Migration() error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user User) (*User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*User, error) {
	var user *User
	if err := r.db.Model(&User{}).Where("username =?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepository) IsUserExistWithSameEmail(email string) bool {
	return r.isUserExistWithCredential("email", email)
}

func (r *userRepository) isUserExistWithCredential(key, value string) bool {
	var user *User
	if err := r.db.Where(fmt.Sprintf("%s =?", key), value).First(&user).Error; err == nil && user != nil {
		return true
	}

	return false
}

func (r *userRepository) IsUserExistWithSameUsername(username string) bool {
	return r.isUserExistWithCredential("username", username)
}

func (r *userRepository) Migration() error {
	return r.db.AutoMigrate(User{})
}
