package binance

import (
    "log"
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/kecbigmt/ccyutils"
)

// struct to receive Binance format JSON
type BinanceTick struct {
  Symbol string `json:"symbol"`
  PriceChange string `json:"priceChange"`
  PriceChangePercent string `json:"priceChangePercent"`
  WeightedAvgPrice string `json:"weightedAvgPrice"`
  PrevClosePrice string `json:"prevClosePrice"`
  LastPrice string `json:"lastPrice"`
  LastQty string `json:"lastQty"`
  BidPrice string`json:"bidPrice"`
  BidQty string `json:"bidQty"`
  AskPrice string `json:"askPrice"`
  AskQty string `json:"askQty"`
  OpenPrice string `json:"openPrice"`
  HighPrice string `json:"highPrice"`
  LowPrice string `json:"lowPrice"`
  Volume string `json:"volume"`
  QuoteVolume string `json:"quoteVolume"`
  OpenTime int64 `json:"openTime"`
  CloseTime int64 `json:"closeTime"`
  FirstId int `json:"firstId"`
  LastId int `json:"lastId"`
  Count int `json:"count"`
}

// normalize BinanceTick struct
func (bt BinanceTick) Norm(currency_pair string) ccyutils.TickData{
  var td ccyutils.TickData
  td.ExchangeName = "Binance"
  td.CurrencyPair = currency_pair
  td.UnixTimestamp = bt.CloseTime / int64(1000)
  td.TickId = bt.LastId
  bp, _ := strconv.ParseFloat(bt.BidPrice, 32)
  td.BestBid = float32(bp)
  ap, _ := strconv.ParseFloat(bt.AskPrice, 32)
  td.BestAsk = float32(ap)
  bq, _ := strconv.ParseFloat(bt.BidQty, 32)
  td.BestBidSize = float32(bq)
  aq, _ := strconv.ParseFloat(bt.AskQty, 32)
  td.BestAskSize = float32(aq)
  ltp, _ := strconv.ParseFloat(bt.LastPrice, 32)
  td.LastPrice = float32(ltp)
  hp, _ := strconv.ParseFloat(bt.HighPrice, 32)
  td.HighPrice = float32(hp)
  lwp, _ := strconv.ParseFloat(bt.LowPrice, 32)
  td.LowPrice = float32(lwp)
  pcp, _ := strconv.ParseFloat(bt.PriceChangePercent, 32)
  td.PriceChangePercent24h = float32(pcp)
  vl, _ := strconv.ParseFloat(bt.Volume, 32)
  td.Volume = float32(vl)
  return td
}

// get ticker
func Ticker(currency_pair string) ccyutils.TickData{
  t_currency_pair := strings.ToUpper(strings.Replace(currency_pair, "_", "", -1)) // XXX_YYY -> xxxyyy
  url := "https://api.binance.com/api/v1/ticker/24hr?symbol="+t_currency_pair
  req, _ := http.NewRequest("GET", url, nil)
  client := new(http.Client)
  resp, _ := client.Do(req)
  defer resp.Body.Close()

  bytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      log.Fatal(err)
  }
  var bt BinanceTick
  if err := json.Unmarshal(bytes, &bt); err != nil {
      log.Fatal(err)
  }
  return bt.Norm(currency_pair)
}
