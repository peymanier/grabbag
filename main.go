package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://apiv2.nobitex.ir/v2/trades/USDTIRT")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%+v\n", data)
}
