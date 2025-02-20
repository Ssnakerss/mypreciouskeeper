package client

import (
	"context"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
)

// Interface for service which sync data between remote and local storage
type SyncService interface {
	//ListLatest gets all assets created after specific time
	ListLatest(
		ctx context.Context,
		userID int64,
		lastUpdated time.Time,
	) ([]*models.Asset, error)
	//Get asset by ID bcause list return names without body...
	Get(
		ctx context.Context,
		userID int64,
		aid int64,
	) (*models.Asset, error)
	//Update or insert asset record by ID and updated time
	UpdateToLatest(
		ctx context.Context,
		asset *models.Asset,
	) (int, error) //result int -  0 - fail, 1 - iserted, 2 - updated
}

type SyncClient struct {
	remoteAssetService SyncService
	localAssetService  SyncService
	lastUpdate         time.Time
}

func NewSyncClient(
	ctx context.Context,
	remoteAssetService SyncService,
	localAssetService SyncService,
) *SyncClient {
	return &SyncClient{
		remoteAssetService: remoteAssetService,
		localAssetService:  localAssetService,
		lastUpdate:         time.Unix(0, 0), //Initial sync time is in the past
	}
}

// Sync get new record from remote storage
// and update local storage with new record of insert
// Than select from local storage
func (s *SyncClient) Sync(ctx context.Context, userID int64) (int, error) {
	//Sync from remote storage
	remoteList, err := s.remoteAssetService.ListLatest(ctx, userID, s.lastUpdate)
	if err != nil {
		return 0, err
	}
	cnt := 0
	for _, remoteAsset := range remoteList {
		asset, err := s.remoteAssetService.Get(ctx, -1, remoteAsset.ID)
		if err != nil {
			return 0, err
		}
		_, err = s.localAssetService.UpdateToLatest(ctx, asset)
		if err != nil {
			return 0, err
		}
		//Remeber max updated tiem for next sync
		//to select only  new recoreds from remote storage
		if remoteAsset.UpdatedAt.Compare(s.lastUpdate) > 0 {
			s.lastUpdate = remoteAsset.UpdatedAt
		}
		cnt++
	}
	return cnt, nil
}
