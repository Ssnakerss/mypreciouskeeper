package models

type User struct {
	ID       int64
	Email    string
	PassHash string
	JWTToken string
}
