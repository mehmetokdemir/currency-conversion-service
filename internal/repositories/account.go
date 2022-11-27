package repositories

import (
	"errors"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account entity.Account) (*entity.Account, error)
	ListUserAccounts(userId uint) ([]entity.Account, error)
	IsUserHasAccountOnGivenCurrency(userId uint, currencyCode string) bool
	GetUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string) (float64, error)
	UpdateUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string, balance float64) error
	Migration() error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) Migration() error {
	return r.db.AutoMigrate(entity.Account{})
}

func (r *accountRepository) Create(account entity.Account) (*entity.Account, error) {
	if err := r.db.Create(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) ListUserAccounts(userId uint) ([]entity.Account, error) {
	var accounts []entity.Account
	if err := r.db.Where("user_id =?", userId).Preload("User").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) IsUserHasAccountOnGivenCurrency(userId uint, currencyCode string) bool {
	var account *entity.Account
	if err := r.db.Where("user_id =?", userId).Where("currency_code =?", currencyCode).First(&account).Error; err != nil {
		return false
	}

	if account == nil {
		return false
	}

	return true
}

func (r *accountRepository) GetUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string) (float64, error) {
	var account *entity.Account
	if err := r.db.Where("user_id =?", userId).Where("currency_code =?", currencyCode).First(&account).Error; err != nil {
		return 0, err
	}

	if account == nil {
		return 0, errors.New("account not found on given currency")
	}

	return account.Balance, nil
}

func (r *accountRepository) UpdateUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string, balance float64) error {
	return r.db.Model(&entity.Account{}).Where("user_id =?", userId).Where("currency_code =?", currencyCode).Update("balance", balance).Error
}
