package binance

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/kecbigmt/ccyutils"
)

// struct to receive bitFlyer format JSON
type BinanceAvailablePairRaw struct {
  Symbol string `json:"symbol"`
}

// get available currency pairs
func AvailablePair() (output ccyutils.AvailablePairs){
    url := "https://api.binance.com/api/v1/ticker/allPrices"
    req, _ := http.NewRequest("GET", url, nil)
    client := new(http.Client)
    resp, err := client.Do(req)
    if err != nil{
      err = fmt.Errorf(`{
        "error_code":"800",
        "component":"available_pair",
        "service":"binance",
        "message":"[Error]Failed to request",
        "detail":%v
      }`, err)
      log.Fatal(err)
    }
    defer resp.Body.Close()
    bytes, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200{
      err = fmt.Errorf(`{
        "error_code":"801",
        "component":"available_pair",
        "service":"binance",
        "message":"[Error]HTTP Error(%v)",
        "detail":%v
      }`, resp.StatusCode, string(bytes))
      return
    }
    var json_receiver []BinanceAvailablePairRaw
    if err = json.Unmarshal(bytes, &json_receiver); err != nil {
      err = fmt.Errorf(`{
        "error_code":"802",
        "component":"available_pair",
        "service":"binance",
        "message":"[Error]Failed to decode JSON",
        "detail":%v
      }`, err)
      log.Fatal(err)
    }
    var sl []string
    for _, v := range json_receiver{
      sl = append(sl, v.Symbol)
    }
    output = ccyutils.AvailablePairs(sl)
    return
}
