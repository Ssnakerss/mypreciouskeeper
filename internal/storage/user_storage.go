package storage

import (
	"context"
	imsql "database/sql"
	"errors"
	"strings"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
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
		ID:       -1,
	}

	sql := ` INSERT INTO  public.mpk_users (u_email, u_pass_hash) VALUES ($1, $2) RETURNING id`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	err = s.DB.QueryRowContext(ctx, sql, uemail, upassHash).Scan(&usr.ID)

	//TO-DO -  check why not working errors.As(err, &pgErr)
	if err != nil {
		// var pgErr *pgconn.PgError
		// if errors.As(err, &pgErr) {
		if strings.Contains(err.Error(), "23505") {
			//если пользователь уже существует, возвращаем ошибку
			// if pgErr.Code == "23505" {
			return usr, apperrs.ErrUserAlreadyExists
			// }
		}
	}
	return usr, err
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
