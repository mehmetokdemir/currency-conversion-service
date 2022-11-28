package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/internal/account"
	"github.com/mehmetokdemir/currency-conversion-service/internal/currency"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"time"
)

type IUserService interface {
	CreateUser(user User) (*User, error)
	CreateToken(username, password string) (*LoginResponse, error)
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

	hashedPassword, err := s.hashPassword(user.Password)
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

	return &user, err
}

func (s *userService) CreateToken(username, password string) (*LoginResponse, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if ok := s.verifyPassword(user.Password, password); !ok {
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

func (s *userService) verifyPassword(hashedPassword, requestedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestedPassword)); err == nil {
		return true
	}
	return false
}

func (s *userService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
