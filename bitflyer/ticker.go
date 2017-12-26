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
  BestBid float32`json:"best_bid"`
  BestAsk float32 `json:"best_ask"`
  BestBidSize float32 `json:"best_bid_size"`
  BestAskSize float32 `json:"best_ask_size"`
  TotalBidDepth float32 `json:"total_bid_depth"`
  TotalAskDepth float32 `json:"total_ask_depth"`
  Ltp float32 `json:"ltp"`
  Volume float32 `json:"volume"`
  VolumeByProduct float32 `json:"volume_by_product"`
}

// normalize BitFlyerTick struct
func (bft BitFlyerTick) Norm(currency_pair string) ccyutils.TickData{
  var td ccyutils.TickData
  td.ExchangeName = "bitFlyer"
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
  return td
}

// get ticker
func Ticker(currency_pair string) ccyutils.TickData{
    url := "https://api.bitflyer.jp/v1/ticker?product_code="+currency_pair
    req, _ := http.NewRequest("GET", url, nil)
    client := new(http.Client)
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    bytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    var bft BitFlyerTick
    if err := json.Unmarshal(bytes, &bft); err != nil {
        log.Fatal(err)
    }
    return bft.Norm(currency_pair)
}
