package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/peymanier/grabbag/database"
	"github.com/peymanier/grabbag/pgconv"

	"net/http"
)

type NobitexTrade struct {
	Time   int    `json:"time"`
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Type   string `json:"type"`
}

type NobitexResponse struct {
	Status string         `json:"status"`
	Trades []NobitexTrade `json:"trades"`
}

func NobitexUpdate(ctx context.Context, queries *database.Queries) {
	codes := []string{"USDTIRT", "BTCUSDT", "ETHUSDT"}
	for _, code := range codes {
		if err := NobitexUpdateAsset(ctx, queries, code); err != nil {
			log.Println(err)
		}
	}
}

func NobitexUpdateAsset(ctx context.Context, queries *database.Queries, code string) error {
	resp, err := http.Get(fmt.Sprintf("https://apiv2.nobitex.ir/v2/trades/%s", code))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	response := NobitexResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	price := pgconv.StringToNumeric(response.Trades[0].Price)

	asset, err := queries.CreateOrUpdateAsset(context.Background(), database.CreateOrUpdateAssetParams{
		Code:      code,
		Price:     price,
		UpdatedAt: pgconv.TimeToTimestamptz(time.Now().UTC()),
	})
	if err != nil {
		return err
	}

	_, err = queries.CreateAssetPriceLog(context.Background(), database.CreateAssetPriceLogParams{
		AssetID: pgconv.Int64ToInt8(asset.ID),
		Price:   price,
	})
	if err != nil {
		return err
	}

	return nil
}
