package account

import (
	"errors"
	"gorm.io/gorm"
)

type IAccountRepository interface {
	CreateAccount(account Account) (*Account, error)
	ListUserAccounts(userId uint) ([]Account, error)
	IsUserHasAccountOnGivenCurrency(userId uint, currencyCode string) bool
	GetUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string) (float64, error)
	UpdateUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string, balance float64) error
	Migration() error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) IAccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) Migration() error {
	return r.db.AutoMigrate(Account{})
}

func (r *accountRepository) CreateAccount(account Account) (*Account, error) {
	if err := r.db.Create(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) ListUserAccounts(userId uint) ([]Account, error) {
	var accounts []Account
	if err := r.db.Where("user_id =?", userId).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) IsUserHasAccountOnGivenCurrency(userId uint, currencyCode string) bool {
	var account *Account
	if err := r.db.Where("user_id =?", userId).Where("currency_code =?", currencyCode).First(&account).Error; err != nil {
		return false
	}

	if account == nil {
		return false
	}

	return true
}

func (r *accountRepository) GetUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string) (float64, error) {
	var account *Account
	if err := r.db.Select("balance").Where("user_id =?", userId).Where("currency_code =?", currencyCode).First(&account).Error; err != nil {
		return 0, err
	}

	if account == nil {
		return 0, errors.New("account not found on given currency")
	}

	return account.Balance, nil
}

func (r *accountRepository) UpdateUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string, balance float64) error {
	return r.db.Model(&Account{}).Where("user_id =?", userId).Where("currency_code =?", currencyCode).Update("balance", balance).Error
}
