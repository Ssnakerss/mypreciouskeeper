package storage

import (
	"context"
	imsql "database/sql"
	"errors"

	"github.com/Ssnakerss/mypreciouskeeper/internal/domain/models"
)

// CreateUSer insert user record into mpk_users table
func (s *DBStorage) CreateUser(pctx context.Context,
	uemail string,
	upassHash string,
) (usr *models.User, err error) {
	usr = &models.User{
		Email:    uemail,
		PassHash: upassHash,
	}

	sql := ` INSERT INTO  public.mpk_users (u_email, u_pass_hash) VALUES ($1, $2) RETURNING id`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	err = s.DB.QueryRowContext(ctx, sql, uemail, upassHash).Scan(&usr.ID)
	if err != nil {
		return nil, err
	}

	return usr, err
}

// GetUser get user record from mpk_users table
func (s DBStorage) GetUser(pctx context.Context, uemail string) (usr *models.User, err error) {
	usr = &models.User{ID: -1}
	sql := ` select id, u_email, u_pass_hash from public.mpk_users where u_email = $1`
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
