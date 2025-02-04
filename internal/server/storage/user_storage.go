package storage

import (
	"context"
	imsql "database/sql"
	"errors"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateUSer insert user record into mpk_users table
func (s *DBStorage) CreateUser(pctx context.Context,
	user *models.User,
) (*models.User, error) {

	sql := ` INSERT INTO  public.mpk_users (u_email, u_pass_hash) VALUES ($1, $2) RETURNING id`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, sql, user.Email, user.PassHash).Scan(&user.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		//Return error when constrain violation -  user already exists
		if pgErr.Code == "23505" {
			return user, apperrs.ErrUserAlreadyExists
		}
	}
	return user, err
}

// GetUser get user record from mpk_users table
func (s DBStorage) GetUser(pctx context.Context, uemail string) (usr *models.User, err error) {
	usr = &models.User{ID: -1}
	sql := ` SELECT id, u_email, u_pass_hash FROM public.mpk_users WHERE u_email = $1`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	err = s.DB.QueryRowContext(ctx, sql, uemail).Scan(&usr.ID, &usr.Email, &usr.PassHash)
	if err != nil {
		if errors.Is(err, imsql.ErrNoRows) {
			return usr, nil
		}
		return nil, err
	}
	return usr, nil
}
