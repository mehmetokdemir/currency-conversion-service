package dto

type AccountWallet struct {
	CurrencyCode string  `json:"currency_code"`
	Balance      float64 `json:"balance"`
}

type UserAccountResponse struct {
	Id       uint            `json:"id"`       // ID of the user
	Username string          `json:"username"` // Username of the user
	Wallets  []AccountWallet `json:"wallets"`  // Wallets of the user
}
