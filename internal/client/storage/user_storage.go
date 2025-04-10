package storage

import (
	"context"
	imsql "database/sql"
	"errors"

	"math/rand"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/mattn/go-sqlite3"
)

// CreateUSer insert user record into mpk_users table
func (s *Storage) CreateUser(pctx context.Context,
	user *models.User,
) (*models.User, error) {

	sql := `
	INSERT INTO  mpk_users 
	(id, u_local_id, u_email, u_pass_hash, u_created_at, u_updated_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6) 
	`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	//Generate local id for user
	//Local ID -  always negative
	//When creating put both ID  and Local ID in user record
	//Remote ID should be updated after user creation and remote login

	userLocalId := rand.Int63() * -1

	_, err := s.db.ExecContext(ctx, sql,
		userLocalId,
		userLocalId,
		user.Email,
		user.PassHash,
		user.CreatedAt.Unix(),
		user.UpdatedAt.Unix(),
	)

	//Checking for errors
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			return user, apperrs.ErrUserAlreadyExists
		}
	}
	return user, err
}

// GetUser get user record from mpk_users table
func (s Storage) GetUser(pctx context.Context, uemail string) (usr *models.User, err error) {
	usr = &models.User{ID: -1}
	sql := `
	SELECT id, u_email, u_pass_hash 
	FROM mpk_users 
	WHERE u_email = $1`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	err = s.db.QueryRowContext(ctx, sql, uemail).Scan(&usr.ID, &usr.Email, &usr.PassHash)
	if err != nil {
		if errors.Is(err, imsql.ErrNoRows) {
			return nil, apperrs.ErrUserNotFound
		}
		return nil, err
	}
	return usr, nil
}
