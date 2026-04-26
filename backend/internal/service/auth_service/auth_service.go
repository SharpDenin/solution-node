package auth_service

import (
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
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

func (s *AuthService) Register(ctx context.Context, fullName, login, password, role, position string) error {
	if fullName == "" || login == "" || password == "" || role == "" {
		return errors.New("invalid input")
	}

	if role != "node" && role != "phenophase" && role != "admin" {
		return errors.New("invalid role")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var positionPtr *string
	if strings.TrimSpace(position) != "" {
		v := strings.TrimSpace(position)
		positionPtr = &v
	}

	user := &models.User{
		FullName:     strings.TrimSpace(fullName),
		Login:        strings.TrimSpace(login),
		PasswordHash: string(hash),
		Position:     positionPtr,
	}

	err = s.userRepo.Create(ctx, user, role)
	if err != nil {
		if err.Error() == "user already exists" {
			return err
		}
		if err.Error() == "invalid role" {
			return err
		}
		return errors.New("failed to create user")
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, fullName, login, password string) (string, error) {
	if fullName == "" || login == "" || password == "" {
		return "", errors.New("invalid credentials")
	}

	user, roleName, err := s.userRepo.GetByLogin(ctx, strings.TrimSpace(login))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !sameFullName(user.FullName, fullName) {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.tokenMgr.GenerateToken(user.ID.String(), roleName)
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID string) (*responses.CurrentUserResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, roleName, err := s.userRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &responses.CurrentUserResponse{
		ID:       user.ID.String(),
		FullName: user.FullName,
		Login:    user.Login,
		Role:     roleName,
		Position: user.Position,
	}, nil
}

func sameFullName(dbName, inputName string) bool {
	return strings.EqualFold(
		strings.TrimSpace(dbName),
		strings.TrimSpace(inputName),
	)
}
