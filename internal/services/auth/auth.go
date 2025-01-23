package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/domain/models"
	"github.com/Ssnakerss/mypreciouskeeper/internal/lib"
	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	CreateUser(ctx context.Context, email string, passHash string) (usr *models.User, err error)
	GetUser(ctx context.Context, email string) (usr *models.User, err error)
	Close() error
}

type Auth struct {
	l        *slog.Logger
	u        UserStorage
	tokenTTL time.Duration
}

// New creates a new instance of auth
func New(l *slog.Logger, u UserStorage, tokenTTL time.Duration) *Auth {
	return &Auth{
		l:        l,
		u:        u,
		tokenTTL: tokenTTL,
	}
}

// Register creates new user with email and password
// First check if same email exists - returns error
// if not - create new user, password replaced with hash
func (a Auth) RegisterUser(ctx context.Context,
	email string,
	pass string,
) (usrID int64, err error) {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "Auth.Register"
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
	usr, err := a.u.CreateUser(ctx, email, string(passHash))
	if err != nil {
		l.Error("failed to create new user", logger.Err(err))
		return -1, err
	}

	return usr.ID, nil
}

// Login authorize user by email and password
// first get user from storage by email
// than compare password with hash from storage and return user if correct
func (a Auth) Login(ctx context.Context, email string, pass string) (token string, err error) {
	who := "Auth.Login"
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
