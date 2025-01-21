package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/domain/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	dsn     string
	timeout time.Duration
	db      *sql.DB
}

// New create an instance of DBStorage with dsn and timeout
func New(pctx context.Context, dsn string, timeout time.Duration) (*DBStorage, error) {
	//open connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	//check connection
	ctx, cancel := context.WithTimeout(pctx, timeout)
	defer cancel()
	err = db.PingContext(ctx)

	return &DBStorage{
		dsn:     dsn,
		timeout: timeout,
		db:      db,
	}, nil
}

// Close close active DB connection
func (s *DBStorage) Close() error {
	return s.db.Close()
}

// Prepare database schema - creates tables and indexes if not exists
func (s *DBStorage) Prepare(pctx context.Context) error {
	sql := `
	CREATE TABLE IF NOT EXISTS public.mpk_users
(
    id bigserial NOT NULL ,
    email text COLLATE pg_catalog."default" NOT NULL,
    pass_hash text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT mpk_users_pkey PRIMARY KEY (id),
	CONSTRAINT mpk_users_unq UNIQUE (email)
);
CREATE INDEX IF NOT EXISTS mpk_users_idx ON public.mpk_users (email)
;
CREATE TABLE IF NOT EXISTS public.mpk_assets
(
    id bigserial NOT NULL ,
    user_id bigint NOT NULL ,
    label text COLLATE pg_catalog."default" NOT NULL,
    asset bytea NOT NULL,
    CONSTRAINT mpk_assets_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS mpk_assets_idx ON public.mpk_assets (user_id)
	`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, sql)

	if err != nil {
		return err
	}
	return nil
}

// CreateUSer insert user record into mpk_users table
func (s *DBStorage) CreateUser(pctx context.Context, email string, passHash string) (usr *models.User, err error) {
	sql := ` insert into  public.mpk_users (email, pass_hash) values  ($1, $2)`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	_, err = s.db.ExecContext(ctx, sql, email, passHash)
	if err != nil {
		return nil, err
	}
	usr, err = s.GetUser(pctx, email)
	return usr, err
}

// GetUser get user record from mpk_users table
func (s DBStorage) GetUser(pctx context.Context, email string) (usr *models.User, err error) {
	usr = &models.User{}
	sql := ` select id, email, pass_hash from public.mpk_users where email = $1`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	err = s.db.QueryRowContext(ctx, sql, email).Scan(&usr.ID, &usr.Email, &usr.PassHash)
	if err != nil {
		return nil, err
	}
	return usr, nil

}
