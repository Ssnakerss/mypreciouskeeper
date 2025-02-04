package services

import (
	"context"
	"errors"
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

// Register creates new user with email and password
// First check if same email exists - returns error
// if not - create new user, password replaced with hash
// gRPC mapping  -  Register
func (a AuthService) RegisterUser(ctx context.Context,
	email string,
	pass string,
) (usrID int64, err error) {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AuthService.Register"
	l := a.l.With(slog.String("who", who), slog.String("email", email))
	l.Info("registering new user")

	//check if user already exist
	if _, err := a.u.GetUser(ctx, email); err != nil {
		if errors.Is(err, apperrs.ErrUserAlreadyExists) {
			l.Warn("user already exist")
			return -1, apperrs.ErrUserAlreadyExists
		} else {
			l.Error("error checking user exist", logger.Err(err))
			return -1, err
		}
	}

	//generating hash for password
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		l.Error("failed to generate password hash", logger.Err(err))
	}

	//Saving new user to storage
	user := &models.User{
		Email:    email,
		PassHash: string(passHash),
	}
	usr, err := a.u.CreateUser(ctx, user)
	if err != nil {
		l.Error("failed to create new user", logger.Err(err))
		return -1, err
	}

	return usr.ID, nil
}

// Login authorize user by email and password
// first get user from storage by email
// than compare password with hash from storage and return user if correct
// gRPC mapping  -  Login
func (a AuthService) Login(ctx context.Context, email string, pass string) (token string, err error) {
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
