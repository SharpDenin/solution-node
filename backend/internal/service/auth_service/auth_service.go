package auth_service

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"
	"time"

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

	if user.IsDeleted {
		return "", errors.New("user deleted")
	}

	if user.IsBlocked {
		return "", errors.New("user blocked")
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

	if user.IsDeleted {
		return nil, errors.New("user deleted")
	}

	if user.IsBlocked {
		return nil, errors.New("user blocked")
	}

	return &responses.CurrentUserResponse{
		ID:       user.ID.String(),
		FullName: user.FullName,
		Login:    user.Login,
		Role:     roleName,
		Position: user.Position,
	}, nil
}

func (s *AuthService) GetAllUsers(ctx context.Context) ([]responses.UserAdminResponse, error) {
	users, rolesByUserID, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]responses.UserAdminResponse, 0, len(users))

	for _, user := range users {
		result = append(result, responses.UserAdminResponse{
			ID:        user.ID.String(),
			FullName:  user.FullName,
			Login:     user.Login,
			Role:      rolesByUserID[user.ID],
			Position:  user.Position,
			IsBlocked: user.IsBlocked,
			IsDeleted: user.IsDeleted,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *AuthService) UpdateUser(ctx context.Context, userID string, req requests.UpdateUserRequest) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	existingUser, existingRole, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if isSystemAdmin(existingUser, existingRole) {
		return errors.New("system admin cannot be modified")
	}

	fullName := strings.TrimSpace(req.FullName)
	login := strings.TrimSpace(req.Login)
	role := strings.TrimSpace(req.Role)

	if fullName == "" || login == "" || role == "" {
		return errors.New("invalid input")
	}

	if role != "node" && role != "phenophase" && role != "admin" {
		return errors.New("invalid role")
	}

	var position *string
	if req.Position != nil && strings.TrimSpace(*req.Position) != "" {
		v := strings.TrimSpace(*req.Position)
		position = &v
	}

	user := &models.User{
		ID:       id,
		FullName: fullName,
		Login:    login,
		Position: position,
	}

	return s.userRepo.Update(ctx, user, role)
}

func (s *AuthService) DeleteUser(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	user, roleName, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if isSystemAdmin(user, roleName) {
		return errors.New("system admin cannot be deleted")
	}

	return s.userRepo.SoftDelete(ctx, id)
}

func (s *AuthService) RestoreUser(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	user, roleName, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if isSystemAdmin(user, roleName) {
		return errors.New("system admin cannot be restored")
	}

	return s.userRepo.Restore(ctx, id)
}

func (s *AuthService) BlockUser(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	user, roleName, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if isSystemAdmin(user, roleName) {
		return errors.New("system admin cannot be blocked")
	}

	return s.userRepo.Block(ctx, id)
}

func (s *AuthService) UnblockUser(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	user, roleName, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if isSystemAdmin(user, roleName) {
		return errors.New("system admin cannot be unblocked")
	}

	return s.userRepo.Unblock(ctx, id)
}

func sameFullName(dbName, inputName string) bool {
	return strings.EqualFold(
		strings.TrimSpace(dbName),
		strings.TrimSpace(inputName),
	)
}

func isSystemAdmin(user *models.User, roleName string) bool {
	if user == nil {
		return false
	}

	return strings.EqualFold(strings.TrimSpace(user.Login), "admin") &&
		strings.EqualFold(strings.TrimSpace(roleName), "admin")
}
