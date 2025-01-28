package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	dsn     string
	timeout time.Duration
	DB      *sql.DB
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
		DB:      db,
	}, nil
}

// Close close active DB connection
func (s *DBStorage) Close() error {
	return s.DB.Close()
}

// Prepare database schema - creates tables and indexes if not exists
func (s *DBStorage) Prepare(pctx context.Context) error {
	sql := `
	CREATE TABLE IF NOT EXISTS public.mpk_users
(
    id bigserial NOT NULL ,
    u_email text COLLATE pg_catalog."default" NOT NULL,
    u_pass_hash text COLLATE pg_catalog."default" NOT NULL,
	u_created_at timestamp with time zone NOT NULL DEFAULT now(),
    u_updated_at timestamp with time zone NOT NULL DEFAULT now(),
	
    CONSTRAINT mpk_users_pkey PRIMARY KEY (id),
	CONSTRAINT mpk_users_unq UNIQUE (u_email)
);
CREATE INDEX IF NOT EXISTS mpk_users_idx ON public.mpk_users (u_email)
;
CREATE TABLE IF NOT EXISTS public.mpk_assets
(
    id bigserial NOT NULL ,
    a_user_id bigint NOT NULL ,
	a_type text COLLATE pg_catalog."default" NOT NULL,
    a_sticker text COLLATE pg_catalog."default" NOT NULL,
    a_body bytea NOT NULL,
	a_created_at timestamp with time zone NOT NULL DEFAULT now(),
    a_updated_at timestamp with time zone NOT NULL DEFAULT now(),
	a_deleted_yn "char" NOT NULL DEFAULT 'N'::"char",
    a_deleted_at timestamp with time zone,
	
    CONSTRAINT mpk_assets_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS mpk_assets_idx ON public.mpk_assets (a_user_id)
	`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, sql)

	if err != nil {
		return err
	}
	return nil
}
