package main

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/peymanier/grabbag/database"
	"github.com/peymanier/grabbag/pgconv"
)

type AssetsWithPriceChangesDTO struct {
	ID              int64
	Code            string
	Price           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PercentChange4h *float64
	Change4h        *string
	PercentChange1d *float64
	Change1d        *string
	PercentChange7d *float64
	Change7d        *string
}

func calculateAssetPriceChange(firstPrice, priceChange pgtype.Numeric) *float64 {
	first := pgconv.NumericToFloat64(firstPrice)
	change := pgconv.NumericToFloat64(priceChange)

	if first == nil || change == nil {
		return nil
	}

	if *first == 0 {
		return nil
	}

	res := *change / *first
	return &res

}

func AssetToDTO(asset database.ListAssetsWithPriceChangesRow) AssetsWithPriceChangesDTO {
	percentChange4h := calculateAssetPriceChange(asset.First4h, asset.Change4h)
	percentChange1d := calculateAssetPriceChange(asset.First1d, asset.Change1d)
	percentChange7d := calculateAssetPriceChange(asset.First7d, asset.Change7d)

	return AssetsWithPriceChangesDTO{
		ID:              asset.ID,
		Code:            asset.Code,
		Price:           *pgconv.NumericToString(asset.Price),
		CreatedAt:       *pgconv.TimestamptzToTime(asset.CreatedAt),
		UpdatedAt:       *pgconv.TimestamptzToTime(asset.UpdatedAt),
		PercentChange4h: percentChange4h,
		Change4h:        pgconv.NumericToString(asset.Change4h),
		PercentChange1d: percentChange1d,
		Change1d:        pgconv.NumericToString(asset.Change1d),
		PercentChange7d: percentChange7d,
		Change7d:        pgconv.NumericToString(asset.Change7d),
	}
}

func AssetsToDTO(assets []database.ListAssetsWithPriceChangesRow) []AssetsWithPriceChangesDTO {
	var res []AssetsWithPriceChangesDTO
	for _, asset := range assets {
		res = append(res, AssetToDTO(asset))
	}

	return res
}
