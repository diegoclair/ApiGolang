package coinmarketcap

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "os"
)

func main() {
  client := &http.Client{}
  req, err := http.NewRequest("GET","https://pro-api.coinmarketcap.com/v1/cryptocurrency/map", nil)
  if err != nil {
    log.Print(err)
    os.Exit(1)
  }

  q := url.Values{}
  q.Add("symbol", "BTC")

  req.Header.Set("Accepts", "application/json")
  req.Header.Add("X-CMC_PRO_API_KEY", "b266000d-ca86-4eb5-9848-6ac6db75a549")
  req.URL.RawQuery = q.Encode()


  resp, err := client.Do(req);
  if err != nil {
    fmt.Println("Error sending request to server")
    os.Exit(1)
  }
  fmt.Println(resp.Status);
  respBody, _ := ioutil.ReadAll(resp.Body)
  fmt.Println(string(respBody));

}