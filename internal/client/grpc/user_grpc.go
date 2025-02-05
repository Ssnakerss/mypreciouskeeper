package grpcClient

import (
	"context"

	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
)

// Login to remote server with email and password and receive auth token
func (c *GRPCClient) Login(
	ctx context.Context,
	email string,
	pass string) (token string, err error) {
	loginResp, err := c.AuthClient.Login(
		context.Background(),
		&grpcserver.LoginRequest{
			Email: email,
			Pass:  pass,
		},
	)
	if err != nil {
		return "", err
	}
	c.token = loginResp.Token
	return loginResp.Token, nil
}

// Register to remote server with email and password and receive userid
func (c *GRPCClient) Register(
	ctx context.Context,
	email string,
	pass string) (userid int64, err error) {
	registerResp, err := c.AuthClient.Register(
		context.Background(),
		&grpcserver.RegisterRequest{
			Email: email,
			Pass:  pass,
		},
	)

	if err != nil {
		return -1, err
	}
	return registerResp.UserId, nil
}
