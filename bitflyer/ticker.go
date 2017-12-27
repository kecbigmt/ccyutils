package bitflyer

import (
    "net/http"
    "io/ioutil"
    "log"
    "time"
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
    url := "https://api.bitflyer.jp/v1/ticker?product_code="+currency_pair
    req, _ := http.NewRequest("GET", url, nil)
    client := new(http.Client)
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil{
      return
    }

    bytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    var bft BitFlyerTick
    if err := json.Unmarshal(bytes, &bft); err != nil {
        log.Fatal(err)
    }
    tick = bft.Norm(currency_pair)
    return
}
