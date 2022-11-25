package entity

import (
	"gorm.io/gorm"
	"time"
)

type Offer struct {
	Id               uint           `gorm:"primaryKey;autoIncrement"`
	FromCurrencyCode string         `gorm:"not null" binding:"required"`
	ToCurrencyCode   string         `gorm:"not null" binding:"required"`
	ExchangeRate     float64        `gorm:"not null" binding:"required"`
	ExpiresAt        int64          `gorm:"not null" binding:"required"`
	UserId           uint           `gorm:"not null" binding:"required"`
	CreatedAt        time.Time      `json:"created_at,omitempty"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
