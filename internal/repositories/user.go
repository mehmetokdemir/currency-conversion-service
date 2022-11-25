package repositories

import (
	"fmt"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"gorm.io/gorm"
)

// DB QUERIES

type UserRepository interface {
	Create(user entity.User) (*entity.User, error)
	Get(id uint) (*entity.User, error)
	IsUserExistWithSameUsername(username string) bool
	IsUserExistWithSameEmail(email string) bool
	GetUserByUsername(username string) (*entity.User, error)
	Migration() error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user entity.User) (*entity.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Get(id uint) (*entity.User, error) {
	var user *entity.User
	if err := r.db.Where("id =?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user *entity.User
	if err := r.db.Where("username =?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) IsUserExistWithSameEmail(email string) bool {
	return r.isUserExistWithCredential("email", email)
}

func (r *userRepository) isUserExistWithCredential(key, value string) bool {
	var user *entity.User
	if err := r.db.Where(fmt.Sprintf("%s =?", key), value).First(&user).Error; err == nil && user != nil {
		return true
	}

	return false
}

func (r *userRepository) IsUserExistWithSameUsername(username string) bool {
	return r.isUserExistWithCredential("username", username)
}

func (r *userRepository) Migration() error {
	return r.db.AutoMigrate(entity.User{})
}
