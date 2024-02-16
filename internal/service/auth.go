package service

import (
	"context"
	"database/sql"
	"errors"
	"example1/config"
	repository "example1/internal/repository/sqlc/generate"
	"example1/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const (
	productWorker = iota
	warehouseWorker
	admin
)

var (
	TakeClaimsErr             = errors.New("error get user claims from token")
	TokenTimeOutErr           = errors.New("token timeout")
	UserNotFoundErr           = errors.New("user not found")
	WrongLoginOrPasswordErr   = errors.New("wrong login or password")
	FailedGenerateTokenErr    = errors.New("failed to generate token")
	FailedGeneratePasswordErr = errors.New("failed to generate password hash")
)

type authService struct {
	Logger         logger.Logger
	Config         config.Config
	UserRepository repository.Queries
}

func NewAuthService(r repository.Queries, c config.Config) Auth {
	return &authService{logger.Get(), c, r}
}

type Auth interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	RegisterNewUser(ctx context.Context, login string, password string, role string) (userID string, err error)
	GetRole(ctx context.Context, userID string) (string, error)
	Authorize(ctx context.Context, access string) (*AuthInfo, error)
	GetAccessByRefresh(ctx context.Context, refresh string) (string, error)
}

func (s *authService) GetAccessByRefresh(ctx context.Context, refresh string) (string, error) {
	s.Logger.Info("starting service GetAccessByRefresh")

	claims, err := ParseSubject(refresh, s.Config.SecretKey)
	if err != nil {
		s.Logger.Error(TokenTimeOutErr)
		return "", TokenTimeOutErr
	}

	id := claims["id"].(string)
	login := claims["login"].(string)
	role := int32(claims["role"].(float64))

	newAccess, err := NewToken(
		repository.User{
			ID:    id,
			Login: login,
		},
		role,
		time.Duration(s.Config.TTLAccessToken)*time.Minute,
		s.Config.SecretKey)

	if err != nil {
		s.Logger.Error(err)
		return "", err
	}
	err = s.UserRepository.UpdateAccessToken(ctx, repository.UpdateAccessTokenParams{
		AccessToken: newAccess, UserID: id,
	})
	if err != nil {
		s.Logger.Error(err)
		return "", err
	}

	return newAccess, nil
}

type AuthInfo struct {
	ID          string
	Login       string
	Role        int32
	AccessToken string
}

func (s *authService) Authorize(ctx context.Context, access string) (*AuthInfo, error) {
	s.Logger.Info("starting service Authorize")
	claims, err := ParseSubject(access, s.Config.SecretKey)
	if err != nil {
		s.Logger.Error(err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenTimeOutErr
		}
		return nil, err
	}

	id := claims["id"].(string)
	login := claims["login"].(string)
	role := int32(claims["role"].(float64))

	return &AuthInfo{
		ID:          id,
		Login:       login,
		Role:        role,
		AccessToken: access,
	}, nil
}

func (s *authService) Login(ctx context.Context, login string, password string) (string, string, error) {
	s.Logger.Info("starting service Login")

	user, err := s.UserRepository.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.Error(UserNotFoundErr)
			return "", "", UserNotFoundErr
		}

		s.Logger.Error(err)

		return "", "", err
	}

	err = ComparePassword([]byte(password), []byte(user.PasswordHash))
	if err != nil {
		s.Logger.Error(WrongLoginOrPasswordErr, ":", err)
		return "", "", WrongLoginOrPasswordErr
	}

	role, err := s.UserRepository.GetRole(ctx, user.ID)
	if err != nil {
		s.Logger.Error(WrongLoginOrPasswordErr, ":", err)
		return "", "", err
	}

	accessToken, err := NewToken(user, role, time.Duration(s.Config.TTLAccessToken)*time.Minute, s.Config.SecretKey)
	if err != nil {
		s.Logger.Error("failed to generate token")
		return "", "", FailedGenerateTokenErr
	}

	refreshToken, err := NewToken(user, role, time.Duration(s.Config.TTLRefreshToken)*time.Minute, s.Config.SecretKey)
	if err != nil {
		s.Logger.Error("failed to generate token")
		return "", "", FailedGenerateTokenErr
	}

	_ = s.UserRepository.CreateTokens(ctx, repository.CreateTokensParams{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	err = s.UserRepository.UpdateTokens(ctx, repository.UpdateTokensParams{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	if err != nil {
		return "", "", FailedGenerateTokenErr
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RegisterNewUser(ctx context.Context, login string, password string, role string) (userID string, err error) {
	s.Logger.Info("starting service RegisterNewUser")

	passHash, err := HashPassword([]byte(password))
	if err != nil {
		s.Logger.Error(FailedGeneratePasswordErr)

		return "", FailedGeneratePasswordErr
	}

	id := uuid.New()
	err = s.UserRepository.CreateUser(ctx, repository.CreateUserParams{
		ID:           id.String(),
		Login:        login,
		PasswordHash: string(passHash),
	})

	if err != nil {
		s.Logger.Error("failed to save user")
		return "", err
	}

	switch role {
	case "product worker":
		err = s.UserRepository.AddRole(ctx, repository.AddRoleParams{UserID: id.String(), Role: productWorker})
	case "warehouse worker":
		err = s.UserRepository.AddRole(ctx, repository.AddRoleParams{UserID: id.String(), Role: warehouseWorker})
	case "admin":
		err = s.UserRepository.AddRole(ctx, repository.AddRoleParams{UserID: id.String(), Role: admin})
	default:
		err = s.UserRepository.AddRole(ctx, repository.AddRoleParams{UserID: id.String(), Role: productWorker})
	}

	if err != nil {
		s.Logger.Error(err)
		return "", err
	}

	return id.String(), nil
}

func (s *authService) GetRole(ctx context.Context, userID string) (string, error) {
	s.Logger.Info("starting service GetRole")

	role, err := s.UserRepository.GetRole(ctx, userID)
	if err != nil {
		s.Logger.Error(err)
		return "", err
	}

	res := ""
	switch role {
	case 0:
		res = "product worker"
	case 1:
		res = "warehouse worker"
	case 2:
		res = "admin"
	}

	return res, nil
}
