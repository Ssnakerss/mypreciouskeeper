package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
)

// CreaeatAsset insert asset record into mpk_assets table
// Convert data to hexade string before insert
func (s *Storage) CreateAsset(pctx context.Context,
	asset *models.Asset,
) (*models.Asset, error) {

	query := `INSERT INTO mpk_assets (
		id, 
		a_user_id, 
		a_type, 
		a_sticker, 
		a_body,
		a_created_at,
		a_updated_at,
		a_deleted_yn,
		a_deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query,
		asset.ID,
		asset.UserID,
		asset.Type,
		asset.Sticker,
		asset.Body,
		asset.CreatedAt.Unix(),
		asset.UpdatedAt.Unix(),
		asset.DeletedYN,
		asset.DeletedAt.Unix(),
	)
	return asset, err
}

// UpdateAsset update asset record into mpk_assets table by id and user_id
func (s *Storage) UpdateAsset(ctx context.Context, asset *models.Asset) (err error) {
	query := `UPDATE mpk_assets SET 
			a_type = $1,
			a_sticker = $2,
			a_body = $3,
			a_updated_at = $4,
			a_deleted_yn = $5,
			a_deleted_at = $6
		WHERE 
			a_user_id = $7 
			AND id = $8
			`
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query,
		asset.Type,
		asset.Sticker,
		asset.Body,
		asset.UpdatedAt.Unix(),
		asset.DeletedYN,
		asset.DeletedAt.Unix(),
		asset.UserID,
		asset.ID,
	)

	if err != nil {
		return err
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrs.ErrAssetNotFound
	}
	return err
}

// GetAsset get asset record from mpk_assets table
func (s Storage) GetAsset(
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
				a_body, 
				a_created_at,
				a_updated_at,
				a_deleted_yn,
				a_deleted_at
			FROM mpk_assets 
			WHERE 
				a_user_id = $1 
				AND id = $2
			`
	ctx, cancel := context.WithTimeout(pctx, s.timeout)
	defer cancel()
	var createdAt, updatedAt, deletedAt int64

	err = s.db.
		QueryRowContext(ctx, query, userID, assetID).
		Scan(
			&asset.Type,
			&asset.Sticker,
			&asset.Body,
			&createdAt,
			&updatedAt,
			&asset.DeletedYN,
			&deletedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrs.ErrAssetNotFound
		}
		return nil, err
	}

	asset.CreatedAt = time.Unix(createdAt, 0)
	asset.UpdatedAt = time.Unix(updatedAt, 0)
	asset.DeletedAt = time.Unix(deletedAt, 0)
	return asset, err

	// return lib.Decompress(data)
}

// DeleteAsset update mpk_assets table by user_id and asset_id set a_deleted_yn = 'Y' and a_deleted_at = now()
func (s *Storage) DeleteAsset(ctx context.Context, userID int64, assetID int64) (err error) {
	query := `
	UPDATE mpk_assets SET
		 a_deleted_yn = 'Y',
		 a_deleted_at = $1
	WHERE
		a_user_id = $2
		AND id = $3
		`
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	_, err = s.db.ExecContext(ctx, query, time.Now().Unix(), userID, assetID)
	return err

}

// ListAssets list assets from mpk_assets table by user_id,  type and sticker 'LIKE' condition
func (s *Storage) ListAssets(ctx context.Context,
	userID int64,
	atype string,
	asticker string) (assets []*models.Asset, err error) {
	var params []any
	params = append(params, userID)
	if atype != "" {
		params = append(params, atype)
	}
	if asticker != "" {
		params = append(params, "%"+asticker+"%")
	}

	query := `SELECT
				id,
				a_type,
				a_sticker,
				a_body,
				a_created_at,
				a_updated_at
			FROM mpk_assets
			WHERE
				a_user_id = $1
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

	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdAt, updatedAt int64
	for rows.Next() {
		asset := &models.Asset{
			UserID: userID,
		}
		err = rows.Scan(&asset.ID, &asset.Type, &asset.Sticker, &asset.Body, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		asset.CreatedAt = time.Unix(createdAt, 0)
		asset.UpdatedAt = time.Unix(updatedAt, 0)

		assets = append(assets, asset)
	}
	return assets, nil
}
