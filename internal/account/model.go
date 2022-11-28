package account

import (
	"gorm.io/gorm"
	"time"
)

// Account Gorm model
type Account struct {
	CurrencyCode string `gorm:"primaryKey;autoIncrement:false"`
	UserId       uint   `gorm:"primaryKey;autoIncrement:false"`
	Balance      float64
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// WalletAccount http response
type WalletAccount struct {
	CurrencyCode string  `json:"currency_code"`
	Balance      float64 `json:"balance"`
}
