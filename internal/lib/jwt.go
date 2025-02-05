package lib

import (
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/golang-jwt/jwt"
)

const (
	// TODO - get app secret from config
	AppSecret   = "poor secret"
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
