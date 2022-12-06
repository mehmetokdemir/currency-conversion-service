package user

import (
	// Go imports
	"errors"
	"os"
	"strings"
	"time"

	// External imports
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	// Internal imports
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
)

type IUserService interface {
	CreateUser(user User) (*User, error)
	CreateToken(username, password string) (*LoginResponse, error)
	VerifyPassword(hashedPassword, requestedPassword string) bool
	HashPassword(password string) (string, error)
}

type userService struct {
	config          config.Config
	userRepository  IUserRepository
	accountService  account.IAccountService
	currencyService currency.Service
}

func NewUserService(userRepository IUserRepository, config config.Config, currencyService currency.Service, accountService account.IAccountService) IUserService {
	return &userService{userRepository: userRepository, config: config, currencyService: currencyService, accountService: accountService}
}

func (s *userService) CreateUser(user User) (*User, error) {

	if s.userRepository.IsUserExistWithSameEmail(user.Email) || s.userRepository.IsUserExistWithSameUsername(user.Username) {
		return nil, errors.New("duplicated users")
	}

	if ok := s.currencyService.CheckIsCurrencyCodeExist(strings.ToUpper(user.DefaultCurrencyCode)); !ok {
		return nil, errors.New("currency not found")
	}

	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	createdUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// First create default account for given currency
	if _, err = s.accountService.CreateUserAccount(createdUser.Id, user.DefaultCurrencyCode, true); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) CreateToken(username, password string) (*LoginResponse, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if ok := s.VerifyPassword(user.Password, password); !ok {
		return nil, errors.New("password mismatch")
	}

	tk := &dto.Token{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(50 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		return nil, errors.New("can not sign jwt")
	}

	return &LoginResponse{
		RegisterResponse: RegisterResponse{
			Username: user.Username,
			Email:    user.Email,
		},
		TokenHash: tokenString,
	}, nil
}

func (s *userService) VerifyPassword(hashedPassword, requestedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestedPassword)); err == nil {
		return true
	}
	return false
}

func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
