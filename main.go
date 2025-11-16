package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/peymanier/grabbag/database"

	"log"
	"net/http"
	"os"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := database.New(conn)

	resp, err := http.Get("https://apiv2.nobitex.ir/v2/trades/USDTIRT")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	response := NobitexResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return
	}

	price, err := strconv.Atoi(response.Trades[0].Price)
	if err != nil {
		log.Println(err)
		return
	}

	priceBigInt := big.NewInt(int64(price))
	asset, err := queries.CreateAsset(context.Background(), database.CreateAssetParams{
		Code:  "USDT",
		Name:  "USD Tether",
		Price: pgtype.Numeric{Int: priceBigInt, Valid: true},
	})
	if err != nil {
		log.Println(err)
		return
	}

	assetPriceLog, err := queries.CreateAssetPriceLog(context.Background(), database.CreateAssetPriceLogParams{
		AssetID: pgtype.Int8{Int64: asset.ID, Valid: true},
		Price:   pgtype.Numeric{Int: priceBigInt, Valid: true},
	})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(asset)
	fmt.Println(assetPriceLog)
}
