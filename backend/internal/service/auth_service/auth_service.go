package auth_service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
	GenerateToken(userID, role string) (string, error)
}

type AuthService struct {
	userRepo repository.UserRepository
	tokenMgr TokenManager
}

func NewAuthService(userRepo repository.UserRepository, tokenMgr TokenManager) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		tokenMgr: tokenMgr,
	}
}

func (s *AuthService) Register(ctx context.Context, fullName, login, password string) error {
	if fullName == "" || login == "" || password == "" {
		return errors.New("invalid input")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		FullName:     fullName,
		Login:        login,
		PasswordHash: string(hash),
		Role:         "worker",
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		if err.Error() == "user already exists" {
			return err
		}
		return errors.New("failed to create user")
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	if login == "" || password == "" {
		return "", errors.New("invalid credentials")
	}

	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.tokenMgr.GenerateToken(user.ID.String(), user.Role)
}
