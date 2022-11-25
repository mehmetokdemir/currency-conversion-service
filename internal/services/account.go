package services

import (
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/mehmetokdemir/currency-conversion-service/internal/repositories"
	"strings"
	"time"
)

type AccountService interface {
	CreateUserAccountOnRegistration(userId uint, currencyCode string) error
	CreateUserAccount(userId uint, currencyCode string) error
	ListUserAccounts(userId uint) ([]dto.AccountWallet, error)
	IsUserHasAccountOnGivenCurrency(userId uint, currencyCode string) bool
}

type accountService struct {
	config      config.Config
	accountRepo repositories.AccountRepository
}

func NewAccountService(accountRepository repositories.AccountRepository, config config.Config) AccountService {
	return &accountService{accountRepo: accountRepository, config: config}
}

func (s *accountService) CreateUserAccountOnRegistration(userId uint, currencyCode string) error {
	return s.createAccount(userId, currencyCode, true)
}

func (s *accountService) createAccount(userId uint, currencyCode string, isOnRegistration bool) error {
	var balance float64 = 0
	if isOnRegistration {
		balance = 10000
	}

	acc := entity.Account{
		CurrencyCode: strings.ToUpper(currencyCode),
		UserId:       userId,
		Balance:      balance,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if _, err := s.accountRepo.Create(acc); err != nil {
		return err
	}

	return nil
}

func (s *accountService) CreateUserAccount(userId uint, currencyCode string) error {
	return s.createAccount(userId, currencyCode, false)
}

func (s *accountService) IsUserHasAccountOnGivenCurrency(userId uint, currencyCode string) bool {
	return s.accountRepo.IsUserHasAccountOnGivenCurrency(userId, strings.ToUpper(currencyCode))
}

func (s *accountService) ListUserAccounts(userId uint) ([]dto.AccountWallet, error) {
	accounts, err := s.accountRepo.ListUserAccounts(userId)
	if err != nil {
		return nil, err
	}

	var respondAccounts []dto.AccountWallet
	for _, account := range accounts {
		respondAccounts = append(respondAccounts, dto.AccountWallet{
			CurrencyCode: account.CurrencyCode,
			Balance:      account.Balance,
		})
	}

	return respondAccounts, err
}
