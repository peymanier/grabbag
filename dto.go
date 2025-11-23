package main

import (
	"time"

	"github.com/peymanier/grabbag/database"
)

type AssetDTO struct {
	ID        int64
	Code      string
	Name      string
	Price     string
	CreatedAt time.Time
}

func AssetToDTO(asset database.Asset) AssetDTO {
	return AssetDTO{
		ID:        asset.ID,
		Code:      asset.Code,
		Name:      asset.Name,
		Price:     *NumericToString(asset.Price),
		CreatedAt: *TimestamptzToTime(asset.CreatedAt),
	}
}

func AssetsToDTO(assets []database.Asset) []AssetDTO {
	var res []AssetDTO
	for _, asset := range assets {
		res = append(res, AssetToDTO(asset))
	}

	return res
}
