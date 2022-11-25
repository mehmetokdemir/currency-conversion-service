package entity

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	CurrencyCode string `gorm:"primaryKey;autoIncrement:false"`
	UserId       uint   `gorm:"primaryKey;autoIncrement:false"`
	User         User   `binding:"required" gorm:"foreignKey:UserId"`
	Balance      float64
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
