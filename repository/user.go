package repository

import (
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"gorm.io/gorm"
)

// DB QUERIES

type UserRepository interface {
	Create(user entity.User) (*entity.User, error)
	Get(id uint) (*entity.User, error)
	CheckUserIsExistWithSameUsername(username string) bool
	CheckUserIsExistWithSameEmail(email string) bool
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
	user := &entity.User{ID: id}
	if err := r.db.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*entity.User, error) {
	user := &entity.User{Username: username}
	if err := r.db.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CheckUserIsExistWithSameEmail(email string) bool {
	user := &entity.User{Email: email}
	if err := r.db.First(user).Error; err != nil {
		return true
	}

	if user != nil {
		return false
	}

	return true
}

func (r *userRepository) CheckUserIsExistWithSameUsername(username string) bool {
	user := &entity.User{Username: username}
	if err := r.db.First(user).Error; err != nil {
		return true
	}

	if user != nil {
		return false
	}

	return true
}

func (r *userRepository) Migration() error {
	return r.db.AutoMigrate(entity.User{})
}
