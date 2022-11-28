package exchange

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestExchangeRepository_CreateOffer(t *testing.T) {
	t.Parallel()
	db, mock := config.ConnectMockDb()
	r := NewExchangeRepository(db)
	o := Offer{
		Id:               uint(uuid.New().ID()),
		FromCurrencyCode: "USD",
		ToCurrencyCode:   "TRY",
		ExchangeRate:     18.63,
		ExpiresAt:        1669602327,
		UserId:           1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(` INSERT INTO "offers" ("from_currency_code","to_currency_code","exchange_rate","expires_at","user_id","created_at","updated_at","deleted_at","id") 
 							VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
		WithArgs(o.FromCurrencyCode, o.ToCurrencyCode, o.ExchangeRate, o.ExpiresAt, o.UserId, o.CreatedAt, o.UpdatedAt, nil, o.Id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "from_currency_code", "to_currency_code", "exchange_rate", "created_at", "updated_at"}).
				AddRow(o.Id, o.FromCurrencyCode, o.ToCurrencyCode, o.ExchangeRate, o.CreatedAt, o.UpdatedAt))

	mock.ExpectCommit()

	offer, err := r.CreateOffer(o)
	assert.Nil(t, err)
	assert.Equal(t, offer.Id, o.Id)
	assert.Equal(t, offer.FromCurrencyCode, o.FromCurrencyCode)
	assert.Equal(t, offer.ToCurrencyCode, o.ToCurrencyCode)
	assert.WithinDuration(t, offer.CreatedAt, o.CreatedAt, 0)
	assert.WithinDuration(t, offer.UpdatedAt, o.UpdatedAt, 0)
	errExpectations := mock.ExpectationsWereMet()
	assert.Nil(t, errExpectations)
}

func TestExchangeRepository_GetExchangeRate(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewExchangeRepository(db)
	expectedExchange := &Exchange{
		FromCurrencyCode: "TRY",
		ToCurrencyCode:   "USD",
		ExchangeRate:     18.63,
		MarkupRate:       0.03,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rows := sqlmock.
		NewRows([]string{"from_currency_code", "to_currency_code", "exchange_rate", "markup_rate", "created_at", "updated_at"}).
		AddRow(expectedExchange.FromCurrencyCode, expectedExchange.ToCurrencyCode,
			expectedExchange.ExchangeRate, expectedExchange.MarkupRate, expectedExchange.CreatedAt, expectedExchange.UpdatedAt)

	sqlSelectOne := `SELECT * FROM "exchanges" WHERE from_currency_code =$1 AND to_currency_code =$2 AND "exchanges"."deleted_at" IS NULL ORDER BY "exchanges"."from_currency_code" LIMIT 1`

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).WithArgs(expectedExchange.FromCurrencyCode, expectedExchange.ToCurrencyCode).WillReturnRows(rows)

	dbExchange, err := r.GetExchangeRate(expectedExchange.FromCurrencyCode, expectedExchange.ToCurrencyCode)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedExchange, dbExchange)
}

func TestExchangeRepository_GetOffer(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewExchangeRepository(db)
	expectedOffer := &Offer{
		Id:               uint(uuid.New().ID()),
		FromCurrencyCode: "USD",
		ToCurrencyCode:   "TRY",
		ExchangeRate:     18.63,
		ExpiresAt:        1669602327,
		UserId:           1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rows := sqlmock.
		NewRows([]string{"id", "from_currency_code", "to_currency_code", "exchange_rate", "expires_at", "user_id", "created_at", "updated_at"}).
		AddRow(expectedOffer.Id, expectedOffer.FromCurrencyCode, expectedOffer.ToCurrencyCode, expectedOffer.ExchangeRate, expectedOffer.ExpiresAt, expectedOffer.UserId, expectedOffer.CreatedAt, expectedOffer.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "offers" WHERE id =$1 AND "offers"."deleted_at" IS NULL ORDER BY "offers"."id" LIMIT 1`)).
		WithArgs(expectedOffer.Id).WillReturnRows(rows)

	dbOffer, err := r.GetOffer(expectedOffer.Id)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedOffer, dbOffer)
}
