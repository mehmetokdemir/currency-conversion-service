package repository

type CurrencyRepository interface {
	ListAllCurrencies() map[string]string
}
