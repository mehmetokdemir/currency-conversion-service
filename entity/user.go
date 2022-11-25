package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                  uint   `gorm:"primaryKey;autoIncrement" `
	Username            string `gorm:"uniqueIndex;not null" binding:"required"`
	Email               string `gorm:"uniqueIndex;not null" binding:"required"`
	Password            string `gorm:"not null" binding:"required"`
	DefaultCurrencyCode string
	CreatedAt           time.Time      `json:"created_at,omitempty"`
	UpdatedAt           time.Time      `json:"updated_at,omitempty"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
