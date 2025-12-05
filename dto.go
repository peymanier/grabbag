package main

import (
	"time"

	"github.com/peymanier/grabbag/database"
	"github.com/peymanier/grabbag/pgconv"
)

type AssetsWithPriceChangesDTO struct {
	ID        int64
	Code      string
	Price     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Change4H  string
	Change1D  string
	Change7D  string
}

func AssetToDTO(asset database.ListAssetsWithPriceChangesRow) AssetsWithPriceChangesDTO {
	return AssetsWithPriceChangesDTO{
		ID:        asset.ID,
		Code:      asset.Code,
		Price:     *pgconv.NumericToString(asset.Price),
		CreatedAt: *pgconv.TimestamptzToTime(asset.CreatedAt),
		UpdatedAt: *pgconv.TimestamptzToTime(asset.UpdatedAt),
		Change4H:  *pgconv.NumericToString(asset.Change4h),
		Change1D:  *pgconv.NumericToString(asset.Change1d),
		Change7D:  *pgconv.NumericToString(asset.Change7d),
	}
}

func AssetsToDTO(assets []database.ListAssetsWithPriceChangesRow) []AssetsWithPriceChangesDTO {
	var res []AssetsWithPriceChangesDTO
	for _, asset := range assets {
		res = append(res, AssetToDTO(asset))
	}

	return res
}
