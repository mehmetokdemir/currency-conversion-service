package exchange

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

// TODO: New Offer repository

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

type OfferRequest struct {
	FromCurrencyCode string `json:"from_currency_code" extensions:"x-order=1" example:"TRY" validate:"required" valid:"required~from_currency_code|invalid"` // From currency code
	ToCurrencyCode   string `json:"to_currency_code" extensions:"x-order=2" example:"EUR" validate:"required" valid:"required~to_currency_code|invalid"`     // To currency code
}

type OfferResponse struct {
	OfferId          uint    `json:"offer_id" extensions:"x-order=1" example:"4"`             // ID of the exchange rate offer
	FromCurrencyCode string  `json:"from_currency_code" extensions:"x-order=2" example:"TRY"` // From currency code
	ToCurrencyCode   string  `json:"to_currency_code" extensions:"x-order=3" example:"EUR"`   // To currency code
	ExchangeRate     float64 `json:"exchange_rate" extensions:"x-order=4" example:"22.00"`    // Exchange rate with markup rate
}

type AcceptOfferRequest struct {
	OfferId uint    `json:"offer_id" extensions:"x-order=1" example:"4" validate:"required" valid:"required~offer_id|invalid"` // ID of the offer
	Amount  float64 `json:"amount" extensions:"x-order=2" example:"100" validate:"required" valid:"required~amount|invalid"`   // ID of the offer
}
