package service

import (
	"errors"
	"time"

	"auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (a *AuthService) Register(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return a.userRepo.Create(user)
}

func (a *AuthService) Login(email, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := a.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) generateJWT(user *domain.User) (string, error) {
	return "jwt-token-placeholder", nil
}
