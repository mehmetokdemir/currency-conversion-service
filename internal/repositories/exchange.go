package repositories

import (
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type ExchangeRepository interface {
	GetExchangeRate(fromCurrency, toCurrency string) (*entity.Exchange, error)
	CreateOffer(offer entity.Offer) (uint, error)
	GetOffer(id uint) (*entity.Offer, error)
	Migration() error
}

type exchangeRepository struct {
	db *gorm.DB
}

func NewExchangeRepository(db *gorm.DB) ExchangeRepository {
	return &exchangeRepository{
		db: db,
	}
}

func (r *exchangeRepository) Migration() error {
	if err := r.db.AutoMigrate(entity.Offer{}); err != nil {
		return err
	}

	if err := r.db.AutoMigrate(entity.Exchange{}); err == nil && r.db.Migrator().HasTable(&entity.Exchange{}) {
		// If exchanges table empty add to default data
		if err = r.db.First(&entity.Exchange{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			var exchanges = []entity.Exchange{
				entity.Exchange{
					FromCurrencyCode: "TRY",
					ToCurrencyCode:   "USD",
					ExchangeRate:     0.054,
					MarkupRate:       0.4,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				entity.Exchange{
					FromCurrencyCode: "USD",
					ToCurrencyCode:   "TRY",
					ExchangeRate:     18.63,
					MarkupRate:       0.3,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				entity.Exchange{
					FromCurrencyCode: "USD",
					ToCurrencyCode:   "EUR",
					ExchangeRate:     0.96,
					MarkupRate:       0.3,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				entity.Exchange{
					FromCurrencyCode: "EUR",
					ToCurrencyCode:   "USD",
					ExchangeRate:     1.04,
					MarkupRate:       0.2,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				entity.Exchange{
					FromCurrencyCode: "TRY",
					ToCurrencyCode:   "EUR",
					ExchangeRate:     0.052,
					MarkupRate:       0.4,
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
				},
				entity.Exchange{
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

func (r *exchangeRepository) GetExchangeRate(fromCurrency, toCurrency string) (*entity.Exchange, error) {
	var exchange *entity.Exchange
	if err := r.db.Where("from_currency_code =?", fromCurrency).Where("to_currency_code =?", toCurrency).First(&exchange).Error; err != nil {
		return nil, err
	}
	return exchange, nil
}

func (r *exchangeRepository) CreateOffer(offer entity.Offer) (uint, error) {
	if err := r.db.Create(&offer).Error; err != nil {
		return 0, err
	}
	return offer.Id, nil
}

func (r *exchangeRepository) GetOffer(id uint) (*entity.Offer, error) {
	var offer *entity.Offer
	if err := r.db.Where("id =?", id).First(&offer).Error; err != nil {
		return nil, err
	}
	return offer, nil
}
