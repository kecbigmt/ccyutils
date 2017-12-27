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
func (bt BinanceTick) Norm(currency_pair string) ccyutils.Tick{
  var td ccyutils.Tick
  td.ServiceName = "binance"
  td.CurrencyPair = currency_pair
  td.UnixTimestamp = bt.CloseTime / int64(1000)
  td.TickId = bt.LastId
  bp, _ := strconv.ParseFloat(bt.BidPrice, 64)
  td.BestBid = bp
  ap, _ := strconv.ParseFloat(bt.AskPrice, 64)
  td.BestAsk = ap
  bq, _ := strconv.ParseFloat(bt.BidQty, 64)
  td.BestBidSize = bq
  aq, _ := strconv.ParseFloat(bt.AskQty, 64)
  td.BestAskSize = aq
  ltp, _ := strconv.ParseFloat(bt.LastPrice, 64)
  td.LastPrice = ltp
  hp, _ := strconv.ParseFloat(bt.HighPrice, 64)
  td.HighPrice = hp
  lwp, _ := strconv.ParseFloat(bt.LowPrice, 64)
  td.LowPrice = lwp
  pcp, _ := strconv.ParseFloat(bt.PriceChangePercent, 64)
  td.PriceChangePercent24h = pcp
  vl, _ := strconv.ParseFloat(bt.Volume, 64)
  td.Volume = vl
  td.Spread = (td.BestAsk - td.BestBid) / td.BestAsk
  return td
}

// get ticker
func Ticker(currency_pair string) (tick ccyutils.Tick, err error){
  t_currency_pair := strings.ToUpper(strings.Replace(currency_pair, "_", "", -1)) // XXX_YYY -> xxxyyy
  url := "https://api.binance.com/api/v1/ticker/24hr?symbol="+t_currency_pair
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
  var bt BinanceTick
  if err := json.Unmarshal(bytes, &bt); err != nil {
      log.Fatal(err)
  }
  tick = bt.Norm(currency_pair)
  return
}
