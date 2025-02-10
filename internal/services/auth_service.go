package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/lib"
	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) (usr *models.User, err error)
	GetUser(ctx context.Context, email string) (usr *models.User, err error)
	Close() error
}

type AuthService struct {
	l        *slog.Logger
	u        UserStorage
	tokenTTL time.Duration
}

// New creates a new instance of auth
func NewAuthService(l *slog.Logger, u UserStorage, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		l:        l,
		u:        u,
		tokenTTL: tokenTTL,
	}
}

// Login authorize user by email and password
// first get user from storage by email
// than compare password with hash from storage and return user if correct
// gRPC mapping  -  Login
func (a *AuthService) Login(
	ctx context.Context,
	email string,
	pass string,
) (token string, err error) {
	who := "AuthService.Login"
	l := a.l.With(slog.String("who", who), slog.String("email", email))
	l.Info("logging in user")

	usr, err := a.u.GetUser(ctx, email)
	if err != nil {
		l.Error("failed to get user", logger.Err(err))
		return "", fmt.Errorf("%s : %w", who, apperrs.ErrInvalidCredentials)
	} else if usr == nil {
		l.Warn("user not found")
		return "", apperrs.ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(usr.PassHash), []byte(pass))
	if err != nil {
		l.Warn("wrong password")
		return "", apperrs.ErrInvalidCredentials
	}

	token, err = lib.NewJWT(usr, a.tokenTTL)
	if err != nil {
		l.Error("failed to generate jwt token", logger.Err(err))
		return token, err
	}

	return token, nil
}

// Register creates new user with email and password
// USer has unique email address
// Uniq is provided by storage engine to prevent duplicate email
// gRPC mapping  -  Register
func (a *AuthService) Register(
	ctx context.Context,
	email string,
	pass string,
) (int64, error) {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AuthService.Register"
	l := a.l.With(slog.String("who", who), slog.String("email", email))
	l.Info("registering new user")

	//generating hash for password
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		l.Error("failed to generate password hash", logger.Err(err))
	}

	//Saving new user to storage
	newUser := &models.User{
		Email:     email,
		PassHash:  string(passHash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	l.Info("creating new user", "user", newUser)
	newUser, err = a.u.CreateUser(ctx, newUser)
	if err != nil {
		l.Error("failed to create new user", logger.Err(err))
		return -1, err
	}

	return newUser.ID, nil
}

// Close undelying storage
func (a *AuthService) Close() {
	a.u.Close()
}
