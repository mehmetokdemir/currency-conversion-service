package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/mehmetokdemir/currency-conversion-service/config"
	"github.com/mehmetokdemir/currency-conversion-service/dto"
	"github.com/mehmetokdemir/currency-conversion-service/entity"
	"github.com/mehmetokdemir/currency-conversion-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type UserService interface {
	CreateUser(user entity.User) (*entity.User, error)
	CreateToken(username, password string) (*dto.LoginResponse, error)
}

type userService struct {
	config   config.Config
	userRepo repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository, config config.Config) UserService {
	return &userService{userRepo: userRepository, config: config}
}

func (s *userService) CreateUser(user entity.User) (*entity.User, error) {

	if s.userRepo.CheckUserIsExistWithSameEmail(user.Email) || s.userRepo.CheckUserIsExistWithSameUsername(user.Username) {
		return nil, errors.New("duplicated users")
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
