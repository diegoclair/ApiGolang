package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Bitcoin struct {
	Status struct {
		Timestamp string `json:"timestamp"`
		ErrorCode int    `json:"error_code"`
	} `json:"status"`
	Data struct {
		BTC struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			LastUpdate string `json:"last_update"`
			Quote      struct {
				USD struct {
					Price float64 `json:"price"`
				} `json:"USD"`
			} `json:"quote"`
		} `json:"BTC"`
	} `json:"data"`
}

var bitcoinPrice float64

//criar func to get price each hour

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func ReloadBiticoinPrice(t time.Time) {
	bitcoinPrice = bitcoinPriceCoinMarketCap()
	fmt.Println(bitcoinPrice)
}

func main() {
	doEvery(6*1000*time.Millisecond, ReloadBiticoinPrice)
}

func GetBitcoinPrice() float64 {
	if bitcoinPrice != 0 {
		return bitcoinPrice
	} else {
		bitcoinPrice = bitcoinPriceCoinMarketCap()
		return bitcoinPrice
	}
}

func bitcoinPriceCoinMarketCap() float64 {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("symbol", "BTC")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "b266000d-ca86-4eb5-9848-6ac6db75a549")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	bitcoin := Bitcoin{}
	json.Unmarshal(respBody, &bitcoin)
	return bitcoin.Data.BTC.Quote.USD.Price
}
