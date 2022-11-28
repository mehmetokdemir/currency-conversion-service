package config

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectMockDb() (*gorm.DB, sqlmock.Sqlmock) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		log.Println("can not connect mock db", err.Error())
		return nil, nil
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDb,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	return db, mock
}
