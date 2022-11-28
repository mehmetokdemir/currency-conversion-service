package account

import (
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"strings"
	"time"
)

type IAccountService interface {
	CreateUserAccount(userId uint, currencyCode string, isOnRegistration bool) (*Account, error)
	ListUserAccounts(userId uint) ([]WalletAccount, error)
	IsUserHasAccountOnGivenCurrency(userId uint, currencyCode string) bool
	GetUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string) (float64, error)
	UpdateUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string, balance float64) error
}

type accountService struct {
	config      config.Config
	accountRepo IAccountRepository
}

func NewAccountService(accountRepository IAccountRepository, config config.Config) IAccountService {
	return &accountService{accountRepo: accountRepository, config: config}
}

func (s *accountService) CreateUserAccount(userId uint, currencyCode string, isOnRegistration bool) (*Account, error) {
	var balance float64 = 0
	if isOnRegistration {
		balance = 10000
	}

	account := Account{
		CurrencyCode: strings.ToUpper(currencyCode),
		UserId:       userId,
		Balance:      balance,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if _, err := s.accountRepo.CreateAccount(account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *accountService) IsUserHasAccountOnGivenCurrency(userId uint, currencyCode string) bool {
	return s.accountRepo.IsUserHasAccountOnGivenCurrency(userId, strings.ToUpper(currencyCode))
}

func (s *accountService) ListUserAccounts(userId uint) ([]WalletAccount, error) {
	accounts, err := s.accountRepo.ListUserAccounts(userId)
	if err != nil {
		return nil, err
	}

	var respondAccounts []WalletAccount
	for _, account := range accounts {
		respondAccounts = append(respondAccounts, WalletAccount{
			CurrencyCode: account.CurrencyCode,
			Balance:      account.Balance,
		})
	}

	return respondAccounts, err
}

func (s *accountService) GetUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string) (float64, error) {
	return s.accountRepo.GetUserBalanceOnGivenCurrencyAccount(userId, currencyCode)
}

func (s *accountService) UpdateUserBalanceOnGivenCurrencyAccount(userId uint, currencyCode string, amount float64) error {
	existingBalance, err := s.GetUserBalanceOnGivenCurrencyAccount(userId, currencyCode)
	if err != nil {
		return err
	}
	existingBalance += amount
	return s.accountRepo.UpdateUserBalanceOnGivenCurrencyAccount(userId, currencyCode, existingBalance)
}
