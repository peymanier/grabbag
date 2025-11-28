package providers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/peymanier/grabbag/database"
	"github.com/peymanier/grabbag/pgconv"
)

type TGJUCoin struct {
	Name  string
	Price string
}

type TGJUCoinResponse struct {
	Coins []TGJUCoin
}

func TGJUUpdateAssets(ctx context.Context, queries *database.Queries) {
	err := TGJUUpdateCoins(ctx, queries)
	if err != nil {
		log.Println(err)
	}
}

func TGJUUpdateCoins(ctx context.Context, queries *database.Queries) error {
	resp, err := http.Get("https://www.tgju.org/coin")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	var response TGJUCoinResponse
	doc.Find(
		"#main > div.container.table-row-style > div > div > div:nth-child(3) > table > tbody",
	).Find("tr").Each(func(i int, selection *goquery.Selection) {
		name := selection.Get(0).Attr[1].Val
		price := selection.Get(0).Attr[5].Val

		response.Coins = append(response.Coins, TGJUCoin{
			Name:  name,
			Price: price,
		})
	})

	for _, coin := range response.Coins {
		code := GetCoinCode(coin)
		price := GetCoinPrice(coin)

		asset, err := queries.CreateOrUpdateAsset(context.Background(), database.CreateOrUpdateAssetParams{
			Code:      code,
			Price:     price,
			UpdatedAt: pgconv.TimeToTimestamptz(time.Now().UTC()),
		})
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = queries.CreateAssetPriceLog(context.Background(), database.CreateAssetPriceLogParams{
			AssetID: pgconv.Int64ToInt8(asset.ID),
			Price:   price,
		})
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func GetCoinCode(coin TGJUCoin) string {
	var code string
	switch coin.Name {
	case "retail_sekee":
		code = "EMAMI/IRT"
	case "retail_sekeb":
		code = "AZADI/IRT"
	case "retail_nim":
		code = "NIM/IRT"
	case "retail_rob":
		code = "ROB/IRT"
	case "retail_gerami":
		code = "GERAMI/IRT"
	}

	return code
}

func GetCoinPrice(coin TGJUCoin) pgtype.Numeric {
	priceStr := strings.ReplaceAll(coin.Price, ",", "")
	price := pgconv.StringToNumeric(priceStr)

	return price
}
