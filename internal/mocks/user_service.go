package mocks

import (
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user entity.User) (*entity.User, error) {
	ret := m.Called(user)
	var r0 *entity.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*entity.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
