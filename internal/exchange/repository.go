package exchange

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type IExchangeRepository interface {
	GetExchangeRate(fromCurrency, toCurrency string) (*Exchange, error)
	CreateOffer(offer Offer) (*Offer, error)
	GetOffer(id uint) (*Offer, error)
	Migration() error
}

type exchangeRepository struct {
	db *gorm.DB
}

func NewExchangeRepository(db *gorm.DB) IExchangeRepository {
	return &exchangeRepository{
		db: db,
	}
}

func (r *exchangeRepository) GetExchangeRate(fromCurrency, toCurrency string) (*Exchange, error) {
	var exchange *Exchange
	if err := r.db.Debug().Where("from_currency_code =?", fromCurrency).Where("to_currency_code =?", toCurrency).First(&exchange).Error; err != nil {
		return nil, err
	}
	return exchange, nil
}

func (r *exchangeRepository) CreateOffer(offer Offer) (*Offer, error) {
	if err := r.db.Debug().Create(&offer).Error; err != nil {
		return nil, err
	}
	return &offer, nil
}

func (r *exchangeRepository) GetOffer(id uint) (*Offer, error) {
	var offer *Offer
	if err := r.db.Debug().Where("id =?", id).First(&offer).Error; err != nil {
		return nil, err
	}
	return offer, nil
}

func (r *exchangeRepository) Migration() error {
	if err := r.db.AutoMigrate(Offer{}); err != nil {
		return err
	}

	if err := r.db.AutoMigrate(Exchange{}); err == nil && r.db.Migrator().HasTable(&Exchange{}) {
		// If exchanges table empty add to default data
		if err = r.db.First(&Exchange{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			var exchanges = []Exchange{
				Exchange{
					FromCurrencyCode: "TRY",
					ToCurrencyCode:   "USD",
					ExchangeRate:     0.054,
					MarkupRate:       0.009,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				Exchange{
					FromCurrencyCode: "USD",
					ToCurrencyCode:   "TRY",
					ExchangeRate:     18.63,
					MarkupRate:       0.3,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				Exchange{
					FromCurrencyCode: "USD",
					ToCurrencyCode:   "EUR",
					ExchangeRate:     0.96,
					MarkupRate:       0.3,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				Exchange{
					FromCurrencyCode: "EUR",
					ToCurrencyCode:   "USD",
					ExchangeRate:     1.04,
					MarkupRate:       0.2,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				Exchange{
					FromCurrencyCode: "TRY",
					ToCurrencyCode:   "EUR",
					ExchangeRate:     0.052,
					MarkupRate:       0.007,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				Exchange{
					FromCurrencyCode: "EUR",
					ToCurrencyCode:   "TRY",
					ExchangeRate:     19.36,
					MarkupRate:       0.2,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
			}
			if err = r.db.Create(exchanges).Error; err != nil {
				return err
			}
		}

	}

	return nil
}
