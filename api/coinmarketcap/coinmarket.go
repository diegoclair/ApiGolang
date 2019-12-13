package coinmarketcap

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


func getCurrentTime() (hr int, mim int) {
	currentTime := time.Now()
	timeStampString := currentTime.Format("2006-01-02 15:04:05")
	layOut := "2006-01-02 15:04:05"
	timeStamp, err := time.Parse(layOut, timeStampString)
	if err != nil {
		fmt.Println(err)
	}

	hr, min, _ := timeStamp.Clock()

	return hr, min
}

func GetBitcoinPrice(old_hr int, old_min int) (hr int, min int, price float64, newHour bool) {

	c_hr, c_min := getCurrentTime()

	if bitcoinPrice != 0 {
		if c_hr == 0 && c_hr != old_hr {
			bitcoinPrice = bitcoinPriceCoinMarketCap()
			return c_hr, c_min, bitcoinPrice, true
		}else{
			if c_hr >= (old_hr + 1) {
				if old_min <= c_min {
					bitcoinPrice = bitcoinPriceCoinMarketCap()
					return c_hr, c_min, bitcoinPrice, true
				}
			}
			return c_hr, c_min, bitcoinPrice, false
		}
		
	} else {
		bitcoinPrice = bitcoinPriceCoinMarketCap()
		return c_hr, c_min, bitcoinPrice, true
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
