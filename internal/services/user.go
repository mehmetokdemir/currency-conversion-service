package services

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/mehmetokdemir/currency-conversion-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"time"
)

type UserService interface {
	CreateUser(user entity.User) (*entity.User, error)
	CreateToken(username, password string) (*dto.LoginResponse, error)
}

type userService struct {
	config          config.Config
	userRepo        repositories.UserRepository
	currencyService CurrencyService
}

func NewUserService(userRepository repositories.UserRepository, config config.Config, currencyService CurrencyService) UserService {
	return &userService{userRepo: userRepository, config: config, currencyService: currencyService}
}

func (s *userService) CreateUser(user entity.User) (*entity.User, error) {

	if s.userRepo.IsUserExistWithSameEmail(user.Email) || s.userRepo.IsUserExistWithSameUsername(user.Username) {
		return nil, errors.New("duplicated users")
	}

	if ok := s.currencyService.CheckCurrencyCodeIsExist(strings.ToUpper(user.DefaultCurrencyCode)); !ok {
		return nil, errors.New("currency not found")
	}

	hashedPassword, err := s.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	return s.userRepo.Create(user)
}

func (s *userService) CreateToken(username, password string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	fmt.Println("username", user.Username)

	if ok := s.verifyPassword(user.Password, password); !ok {
		return nil, errors.New("password mismatch")
	}

	tk := &dto.Token{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(50 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		return nil, errors.New("can not sign jwt")
	}

	return &dto.LoginResponse{
		RegisterResponse: dto.RegisterResponse{
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
