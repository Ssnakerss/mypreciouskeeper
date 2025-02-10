package client

import (
	"context"

	"github.com/Ssnakerss/mypreciouskeeper/internal/lib"
)

// Auth functions

// Login using remote service first
// If success - register user locally
// If fail - try to login locally
// If fail - return error
func (app *ClientApp) Login(
	ctx context.Context,
	login, password string) (string, error) {
	//Keep logn password for  remote login
	//If first login locally -  after  connection reestabloshed try login remotely
	app.login = login
	app.password = password

	//Try login remotely
	token, err := app.remoteAuthService.Login(ctx, login, password)
	if err != nil {
		//Try login locally
		token, err = app.localAuthService.Login(ctx, login, password)
		//switch mode to local
		app.Workmode = LOCAL
		//return empty  token
		return "", err
	} else {
		remoteUsr, err := lib.VerifyJWTPayload(token)
		if err == nil {
			app.RemoteUserID = remoteUsr.ID
			app.UserName = remoteUsr.Email
			app.L.Info("remote user login", "name", app.UserName, "remote id", app.RemoteUserID)
		}
		//Remote login success
		app.Workmode = REMOTE
		//Try register same user locally with same login and password
		app.LocalUsersID, err = app.localAuthService.Register(ctx, login, password)
		//TODO: update local user record with remote user info by NAME/EMAIL
	}
	app.AuthToken = token
	return token, nil
}

// Register user with remote service first
// Then register with local service
func (app *ClientApp) Register(
	ctx context.Context,
	login, password string) (int64, error) {
	//Try register remotely
	app.L.Info("trying remote regsiter")
	remoteUserID, err := app.remoteAuthService.Register(ctx, login, password)
	if err != nil {
		//Try register locally
		app.L.Error("remote regsiter", "error", err)
		app.L.Info("trying local regsiter")
		localUserID, err := app.localAuthService.Register(ctx, login, password)
		if err != nil {
			app.L.Error("local regsiter", "error", err)
			return 0, err
		}
		//Local register success
		app.Workmode = LOCAL
		return localUserID, nil
	}
	app.Workmode = REMOTE
	return remoteUserID, nil
}
