package main

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/peymanier/grabbag/database"
)

type AssetResponse struct {
	ID        int64              `json:"id"`
	Code      string             `json:"code"`
	Name      string             `json:"name"`
	Price     pgtype.Numeric     `json:"price"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

func AssetToResponse(asset database.Asset) AssetResponse {
	return AssetResponse{
		ID:        asset.ID,
		Code:      asset.Code,
		Name:      asset.Name,
		Price:     asset.Price,
		CreatedAt: asset.CreatedAt,
	}
}

func AssetsToResponse(assets []database.Asset) []AssetResponse {
	var assetResponses []AssetResponse
	for _, asset := range assets {
		assetResponses = append(assetResponses, AssetToResponse(asset))
	}

	return assetResponses
}
