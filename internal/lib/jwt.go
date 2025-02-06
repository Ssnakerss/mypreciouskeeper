package lib

import (
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/golang-jwt/jwt"
)

const (
	// TODO - get app secret from config
	AppSecret   = "f3e58332-a779-4b0c-bb82-5f5ee5673228"
	JWTDuration = time.Hour * 24
)

// NewJWT generate new token with authorized user information in claims
// Token lifetime spec is defined by duration parameter
func NewJWT(user *models.User, duration time.Duration) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err = token.SignedString([]byte(AppSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyJWTPayload verify token payload and extract user data from  it
func VerifyJWTPayload(token string) (*models.User, error) {
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(AppSecret), nil
	})

	if err != nil {
		return nil, apperrs.ErrInvalidToken
	}

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperrs.ErrInvalidToken
	}
	//TODO checking for expired token
	return &models.User{
		ID:    int64(claims["id"].(float64)),
		Email: claims["username"].(string),
	}, nil
}
