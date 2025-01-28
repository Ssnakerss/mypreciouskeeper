package storage

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/domain/models"
)

// CreaeatAsset insert asset record into mpk_assets table
// Convert data to hexade string before insert
func (s *DBStorage) CreateAsset(pctx context.Context,
	asset *models.Asset,
) (assetID int64, err error) {

	//TO-DO - benchmark compresssion and decompression
	// abody, err = lib.Compress(abody)

	//Postgres bytea support insert HEX encoded data

	query := `INSERT INTO public.mpk_assets (
		a_user_id, 
		a_type, 
		a_sticker, 
		a_body
		) VALUES ($1, $2, $3, $4) RETURNING id`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()
	err = s.DB.QueryRowContext(ctx, query,
		asset.UserID,
		asset.Type,
		asset.Sticker,
		`\x`+hex.EncodeToString(asset.Body)).Scan(&asset.ID)

	return asset.ID, err
}

// UpdateAsset update asset record into mpk_assets table by id and user_id
func (s *DBStorage) UpdateAsset(ctx context.Context, asset *models.Asset) (err error) {
	query := `UPDATE public.mpk_assets SET 
			a_type = $1,
			a_sticker = $2,
			a_body = $3 
		WHERE 
			a_user_id = $4 
			AND id = $5
			AND a_deleted_yn = 'N'
			`
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	res, err := s.DB.ExecContext(ctx, query,
		asset.Type,
		asset.Sticker,
		`\x`+hex.EncodeToString(asset.Body),
		asset.UserID, asset.ID)

	if err != nil {
		return err
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrs.ErrAssetNotFound
	}
	return err
}

// GetAsset get asset record from mpk_assets table
func (s DBStorage) GetAsset(
	pctx context.Context,
	userID int64,
	assetID int64) (asset *models.Asset, err error) {
	asset = &models.Asset{
		UserID: userID,
		ID:     assetID,
	}

	query := `SELECT 
				a_type,
				a_sticker,   
				a_body 
			FROM public.mpk_assets 
			WHERE 
				a_user_id = $1 
				AND id = $2
				AND a_deleted_yn = 'N'
			`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()

	err = s.DB.QueryRowContext(ctx, query, userID, assetID).Scan(&asset.Type, &asset.Sticker, &asset.Body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrs.ErrAssetNotFound
		}
		return nil, err
	}
	return asset, err
	// return lib.Decompress(data)
}

// ListAssets list assets from mpk_assets table by user_id,  type and sticker 'LIKE' condition
func (s *DBStorage) ListAssets(ctx context.Context,
	userID int64,
	atype string,
	asticker string) (assets []*models.Asset, err error) {
	query := `SELECT 
				id, 
				a_type,
				a_sticker,   
				a_body 
			FROM public.mpk_assets 
			WHERE 
				a_deleted_yn = 'N'
				AND a_user_id = $1
				`
	if atype != "" {
		query += `
				AND a_type = $2
				`
	}
	if asticker != "" && atype != "" {
		query += `
				AND a_sticker LIKE $3
				`
	}
	if asticker != "" && atype == "" {
		query += `
				AND a_sticker LIKE $2
				`
	}

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var params []any
	params = append(params, userID)
	if atype != "" {
		params = append(params, atype)
	}
	if asticker != "" {
		params = append(params, "%"+asticker+"%")
	}

	rows, err := s.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		asset := &models.Asset{
			UserID: userID,
		}
		err = rows.Scan(&asset.ID, &asset.Type, &asset.Sticker, &asset.Body)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

// DeleteAsset update mpk_assets table by user_id and asset_id set a_deleted_yn = 'Y' and a_deleted_at = now()
func (s *DBStorage) DeleteAsset(ctx context.Context, userID int64, assetID int64) (err error) {
	query := `
	UPDATE public.mpk_assets SET
		 a_deleted_yn = 'Y', 
		 a_deleted_at = now()
	WHERE 
		a_user_id = $1 
		AND id = $2
		AND a_deleted_yn = 'N'
		`
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	_, err = s.DB.ExecContext(ctx, query, userID, assetID)
	return err

}
