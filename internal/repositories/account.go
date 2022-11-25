package repositories

import (
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account entity.Account) (*entity.Account, error)
	ListAccounts(userId uint) ([]entity.Account, error)

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
	return nil, nil
}

func (r *accountRepository) ListAccounts(userId uint) ([]entity.Account, error) {
	return nil, nil
}
