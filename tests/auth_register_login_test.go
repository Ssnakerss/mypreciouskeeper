package tests

import (
	"testing"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/lib"
	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ssnakerss/mypreciouskeeper/tests/suite"
	"github.com/brianvoe/gofakeit"
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t) // Создаём Suite

	// TODO: Подготовить данные для тестовых запросов (случайные)
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, 10)

	// TODO: Сделать нужные запросы
	respReg, err := st.AClient.Register(ctx, &grpcserver.RegisterRequest{
		Email: email,
		Pass:  pass,
	})
	t.Log("user id ", respReg.GetUserId())
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AClient.Login(ctx, &grpcserver.LoginRequest{
		Email: email,
		Pass:  pass,
	})
	require.NoError(t, err)
	t.Log(respLogin)

	// TODO: Проверить результаты
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	// Отмечаем время, в которое бы выполнен логин.
	// Это понадобится для проверки TTL токена
	loginTime := time.Now()

	// Парсим и валидируем токен
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(lib.AppSecret), nil
	})
	// Если ключ окажется невалидным, мы получим соответствующую ошибку
	require.NoError(t, err)

	// Преобразуем к типу jwt.MapClaims, в котором мы сохраняли данные
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	t.Log(claims)

	// Проверяем содержимое токена
	assert.Equal(t, respReg.GetUserId(), int64(claims["id"].(float64)))
	assert.Equal(t, email, claims["username"].(string))

	// Проверяем, что TTL токена примерно соответствует нашим ожиданиям.
	assert.InDelta(t, loginTime.Add(time.Hour*24).Unix(), claims["exp"].(float64), 5)

}
