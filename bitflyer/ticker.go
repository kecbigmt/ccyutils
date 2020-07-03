package bitflyer

import (
    "fmt"
    "time"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/kecbigmt/ccyutils"
)

// struct to receive bitFlyer format JSON
type BitFlyerTick struct {
  ProductCode string `json:"product_code"`
  Timestamp string `json:"timestamp"`
  TickId int `json:"tick_id"`
  BestBid float64`json:"best_bid"`
  BestAsk float64 `json:"best_ask"`
  BestBidSize float64 `json:"best_bid_size"`
  BestAskSize float64 `json:"best_ask_size"`
  TotalBidDepth float64 `json:"total_bid_depth"`
  TotalAskDepth float64 `json:"total_ask_depth"`
  Ltp float64 `json:"ltp"`
  Volume float64 `json:"volume"`
  VolumeByProduct float64 `json:"volume_by_product"`
}

// normalize BitFlyerTick struct
func (bft BitFlyerTick) Norm(currency_pair string) ccyutils.Tick{
  var td ccyutils.Tick
  td.ServiceName = "bitflyer"
  td.CurrencyPair = currency_pair
  ut, _ := time.Parse("2006-01-02T15:04:05", bft.Timestamp)
  td.UnixTimestamp = ut.Unix()
  td.TickId = bft.TickId
  td.BestBid = bft.BestBid
  td.BestAsk = bft.BestAsk
  td.BestBidSize = bft.BestBidSize
  td.BestAskSize = bft.BestAskSize
  td.TotalBidDepth = bft.TotalBidDepth
  td.TotalAskDepth = bft.TotalAskDepth
  td.LastPrice = bft.Ltp
  td.Volume = bft.Volume
  td.Spread = (td.BestAsk - td.BestBid) / td.BestAsk
  return td
}

// get ticker
func Ticker(currency_pair string) (tick ccyutils.Tick, err error){
    t_currency_pair := strings.ToUpper(currency_pair) // XXX_YYY
    url := "https://api.bitflyer.jp/v1/ticker?product_code="+t_currency_pair
    req, _ := http.NewRequest("GET", url, nil)
    client := new(http.Client)
    resp, err := client.Do(req)
    if err != nil{
      err = fmt.Errorf(`{
        "error_code":"100",
        "component":"ticker",
        "service":"bitflyer",
        "message":"[Error]Failed to request",
        "detail":%v
      }`, err)
      return
    }
    defer resp.Body.Close()
    bytes, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200{
      err = fmt.Errorf(`{
        "error_code":"101",
        "component":"ticker",
        "service":"bitflyer",
        "message":"[Error]HTTP Error(%v)",
        "detail":%v
      }`, resp.StatusCode, string(bytes))
      return
    }
    var bft BitFlyerTick
    if err = json.Unmarshal(bytes, &bft); err != nil {
      err = fmt.Errorf(`{
        "error_code":"102",
        "component":"ticker",
        "service":"bitflyer",
        "message":"[Error]Failed to decode JSON",
        "detail":%v
      }`, err)
      return
    }
    tick = bft.Norm(currency_pair)
    return
}
