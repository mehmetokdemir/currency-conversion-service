package user

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()
	db, mock := config.ConnectMockDb()
	r := NewUserRepository(db)
	u := User{
		Id:                  uint(uuid.New().ID()),
		Username:            "test-user",
		Email:               "test-user@gmail.com",
		Password:            "1234",
		DefaultCurrencyCode: "TRY",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(` INSERT INTO "users" ("username","email","password","default_currency_code","created_at","updated_at","deleted_at","id") 
 							VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(u.Username, u.Email, u.Password, u.DefaultCurrencyCode, u.CreatedAt, u.UpdatedAt, nil, u.Id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "username", "created_at", "updated_at"}).
				AddRow(u.Id, u.Email, u.Username, u.CreatedAt, u.UpdatedAt))

	mock.ExpectCommit()

	user, err := r.CreateUser(u)
	assert.Nil(t, err)
	assert.Equal(t, user.Id, u.Id)
	assert.Equal(t, user.Username, u.Username)
	assert.Equal(t, user.Email, u.Email)
	assert.WithinDuration(t, user.CreatedAt, u.CreatedAt, 0)
	assert.WithinDuration(t, user.UpdatedAt, u.UpdatedAt, 0)
	errExpectations := mock.ExpectationsWereMet()
	assert.Nil(t, errExpectations)
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewUserRepository(db)
	expectedUser := &User{
		Id:                  uint(uuid.New().ID()),
		Username:            "test",
		Email:               "test@gmail.com",
		Password:            "test",
		DefaultCurrencyCode: "TRY",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "password", "default_currency_code", "created_at", "updated_at"}).
		AddRow(expectedUser.Id, expectedUser.Username, expectedUser.Email, expectedUser.Password, expectedUser.DefaultCurrencyCode, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	sqlSelectOne := `SELECT * FROM "users" WHERE username =$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).WithArgs(expectedUser.Username).WillReturnRows(rows)

	dbUser, err := r.GetUserByUsername(expectedUser.Username)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedUser, dbUser)
}

func TestUserRepository_IsUserExistWithSameEmail(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewUserRepository(db)
	u := &User{
		Id:                  uint(uuid.New().ID()),
		Username:            "test",
		Email:               "test@gmail.com",
		Password:            "test",
		DefaultCurrencyCode: "TRY",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "password", "default_currency_code", "created_at", "updated_at"}).
		AddRow(u.Id, u.Username, u.Email, u.Password, u.DefaultCurrencyCode, u.CreatedAt, u.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email =$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(u.Email).WillReturnRows(rows)

	exist := r.IsUserExistWithSameEmail(u.Email)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, exist, true)
}

func TestUserRepository_Migration(t *testing.T) {
	//db, mock := config.ConnectMockDb()
}

func TestUserRepository_IsUserExistWithSameUsername(t *testing.T) {
	db, mock := config.ConnectMockDb()
	r := NewUserRepository(db)
	u := &User{
		Id:                  uint(uuid.New().ID()),
		Username:            "test",
		Email:               "test@gmail.com",
		Password:            "test",
		DefaultCurrencyCode: "TRY",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "password", "default_currency_code", "created_at", "updated_at"}).
		AddRow(u.Id, u.Username, u.Email, u.Password, u.DefaultCurrencyCode, u.CreatedAt, u.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username =$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(u.Username).WillReturnRows(rows)

	exist := r.IsUserExistWithSameUsername(u.Username)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, exist, true)
}
