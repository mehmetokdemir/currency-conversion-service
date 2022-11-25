package entity

import (
	"gorm.io/gorm"
	"time"
)

type Exchange struct {
	FromCurrencyCode string  `gorm:"primaryKey;autoIncrement:false"`
	ToCurrencyCode   string  `gorm:"primaryKey;autoIncrement:false"`
	ExchangeRate     float64 `gorm:"not null" binding:"required"`
	MarkupRate       float64
	CreatedAt        time.Time      `json:"created_at,omitempty"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
