package account

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	t.Parallel()
	db, mock := config.ConnectMockDb()
	r := NewAccountRepository(db)

	a := Account{
		CurrencyCode: "EUR",
		UserId:       uint(1),
		Balance:      300,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(` INSERT INTO "accounts" ("balance","created_at","updated_at","deleted_at",currency_code,"user_id") 
 							VALUES ($1,$2,$3,$4,$5,$6) RETURNING "currency_code" AND "user_id"`)).
		WithArgs(a.Balance, a.CreatedAt, a.UpdatedAt, nil, a.CurrencyCode, a.UserId).
		WillReturnRows(
			sqlmock.NewRows([]string{"currency_code", "user_id", "balance", "created_at", "updated_at", "deleted_at"}).
				AddRow(a.CurrencyCode, a.UserId, a.Balance, a.CreatedAt, a.UpdatedAt, a.DeletedAt))

	mock.ExpectCommit()

	account, err := r.CreateAccount(a)
	assert.Nil(t, err)
	assert.Equal(t, account.CurrencyCode, a.CurrencyCode)
	assert.Equal(t, account.UserId, a.UserId)
	assert.Equal(t, account.Balance, a.Balance)
	assert.WithinDuration(t, account.CreatedAt, a.CreatedAt, 0)
	assert.WithinDuration(t, account.UpdatedAt, a.UpdatedAt, 0)
	errExpectations := mock.ExpectationsWereMet()
	assert.Nil(t, errExpectations)
}

func TestAccountRepository_GetUserBalanceOnGivenCurrencyAccount(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewAccountRepository(db)
	userId := uint(1)
	currencyCode := "TRY"

	expected := float64(100)

	rows := sqlmock.
		NewRows([]string{"user_id", "currency_code", "balance"}).
		AddRow(userId, currencyCode, expected)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "balance" FROM "accounts" WHERE user_id =$1 AND currency_code =$2 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."currency_code" LIMIT 1`)).
		WithArgs(userId, currencyCode).WillReturnRows(rows)

	actual, err := r.GetUserBalanceOnGivenCurrencyAccount(userId, currencyCode)

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, expected, actual)
}

func TestAccountRepository_IsUserHasAccountOnGivenCurrency(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewAccountRepository(db)
	userId := uint(1)
	currencyCode := "TRY"

	rows := sqlmock.
		NewRows([]string{"user_id", "currency_code"}).
		AddRow(userId, currencyCode)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE user_id =$1 AND currency_code =$2 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."currency_code" LIMIT 1`)).
		WithArgs(userId, currencyCode).WillReturnRows(rows)

	actual := r.IsUserHasAccountOnGivenCurrency(userId, currencyCode)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.True(t, actual)
}

func TestAccountRepository_ListUserAccounts(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewAccountRepository(db)
	userId := uint(1)
	walletAccounts := []Account{
		{
			CurrencyCode: "TRY",
			UserId:       userId,
			Balance:      float64(370),
			CreatedAt:    time.Now().Local().AddDate(0, 0, -50),
			UpdatedAt:    time.Now().Local().AddDate(0, 0, -20),
		},
		{
			CurrencyCode: "USD",
			UserId:       userId,
			Balance:      float64(125),
			CreatedAt:    time.Now().Local().AddDate(0, 0, -20),
			UpdatedAt:    time.Now().Local().AddDate(0, 0, -7),
		},
		{
			CurrencyCode: "USD",
			UserId:       userId,
			Balance:      float64(125),
			CreatedAt:    time.Now().Local(),
			UpdatedAt:    time.Now().Local(),
		},
	}

	rows := sqlmock.
		NewRows([]string{"currency_code", "user_id", "balance", "created_at", "updated_at", "deleted_at"})

	for _, accountWallet := range walletAccounts {
		rows.AddRow(accountWallet.CurrencyCode, accountWallet.UserId, accountWallet.Balance, accountWallet.CreatedAt, accountWallet.UpdatedAt, accountWallet.DeletedAt)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE user_id =$1 AND "accounts"."deleted_at" IS NULL`)).
		WithArgs(userId).WillReturnRows(rows)

	accounts, err := r.ListUserAccounts(userId)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, accounts, walletAccounts)
}

func TestAccountRepository_UpdateUserBalanceOnGivenCurrencyAccount(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewAccountRepository(db)
	userId := uint(1)
	currencyCode := "TRY"
	balance := float64(2500)
	//mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(` UPDATE "accounts" SET balance = $1,"updated_at" = $2
 							WHERE user_id =$3 AND currency_code =$4 AND "accounts.deleted_at" IS NULL `)).
		WithArgs(balance, time.Now(), userId, currencyCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	err := r.UpdateUserBalanceOnGivenCurrencyAccount(userId, currencyCode, balance)
	assert.Nil(t, err)
	errExpectations := mock.ExpectationsWereMet()
	assert.Nil(t, errExpectations)
}
