package dto

type ExchangeRateOfferRequest struct {
	FromCurrencyCode string `json:"from_currency_code" extensions:"x-order=1" example:"TRY" validate:"required" valid:"required~from_currency_code|invalid"` // From currency code
	ToCurrencyCode   string `json:"to_currency_code" extensions:"x-order=2" example:"EUR" validate:"required" valid:"required~to_currency_code|invalid"`     // To currency code
}

type ExchangeRateOfferResponse struct {
	OfferId          uint    `json:"offer_id" extensions:"x-order=1" example:"4"`             // ID of the exchange rate offer
	FromCurrencyCode string  `json:"from_currency_code" extensions:"x-order=2" example:"TRY"` // From currency code
	ToCurrencyCode   string  `json:"to_currency_code" extensions:"x-order=3" example:"EUR"`   // To currency code
	ExchangeRate     float64 `json:"exchange_rate" extensions:"x-order=4" example:"22.00"`    // Exchange rate with markup rate
}
