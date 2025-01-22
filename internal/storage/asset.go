package storage

import (
	"context"
	"encoding/hex"
)

// CreaeatAsset insert asset record into mpk_assets table
// Convert data to hexade string before insert
func (s *DBStorage) CreateAsset(pctx context.Context,
	auserID int64,
	atype string,
	asticker string,
	abody []byte,
) (assetID int64, err error) {

	//TO-DO - compressing as a part of transport process
	// abody, err = lib.Compress(abody)

	encoded := `\x` + hex.EncodeToString(abody)
	sql := `INSERT INTO public.mpk_assets (a_user_id, a_type, a_sticker, a_body) VALUES ($1, $2, $3, $4) RETURNING id`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()
	err = s.DB.QueryRowContext(ctx, sql, auserID, atype, asticker, encoded).Scan(&assetID)

	return assetID, err
}

// GetAsset get asset record from mpk_assets table
func (s DBStorage) GetAsset(pctx context.Context, assetID int64) (data []byte, err error) {
	sql := `select a_body from public.mpk_assets where id = $1`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	err = s.DB.QueryRowContext(ctx, sql, assetID).Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, err
	// return lib.Decompress(data)
}
