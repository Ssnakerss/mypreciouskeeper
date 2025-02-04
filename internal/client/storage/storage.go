package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db      *sql.DB
	timeout time.Duration
}

func New(storagePath string, timeout time.Duration) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db:      db,
		timeout: timeout,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

// Prepare database schema - creates tables and indexes if not exists
func (s *Storage) Prepare(pctx context.Context) string {
	//User table creation
	errString := ""
	sql := `
	CREATE TABLE IF NOT EXISTS "mpk_users" (
	"id" INTEGER NOT NULL PRIMARY KEY ,
	"u_email" TEXT NOT NULL UNIQUE,
	"u_pass_hash" TEXT NOT NULL,
	"u_created_at" INTEGER NOT NULL,
	"u_updated_at" INTEGER NOT NULL
	)`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, sql)
	if err != nil {
		errString += err.Error() + "\n"
	}

	//Create assets table
	sql = `
	CREATE TABLE "mpk_assets" (
	"id" INTEGER NOT NULL PRIMARY KEY,
	"a_user_id" INTEGER NOT NULL,
	"a_type" TEXT NOT NULL,
	"a_sticker" TEXT NOT NULL,
	"a_body" blob NOT NULL,
	"a_created_at" INTEGER,
	"a_updated_at" INTEGER,
	"a_deleted_yn" TEXT,
	"a_deleted_at" INTEGER

	);`
	ctx1, cancel1 := context.WithTimeout(pctx, s.timeout)
	defer cancel1()

	_, err = s.db.ExecContext(ctx1, sql)

	if err != nil {
		errString += err.Error() + "\n"
	}

	return errString
}
