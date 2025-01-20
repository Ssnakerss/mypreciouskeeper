package lib

import (
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/domain/models"
	"github.com/golang-jwt/jwt"
)

// TO-DO - make more efficient
var secret = "secret"

// NewJWT generate new token with authorized user information in claims
// Token lifetime spec is defined by duration parameter
func NewJWT(user *models.User, duration time.Duration) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err = token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
