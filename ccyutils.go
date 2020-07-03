// common
package ccyutils
import (
  "math"
  "strings"
)

// Ticker
type Tick struct {
  ServiceName string
  CurrencyPair string
  UnixTimestamp int64
  TickId int
  BestBid float64
  BestAsk float64
  BestBidSize float64
  BestAskSize float64
  TotalBidDepth float64
  TotalAskDepth float64
  LastPrice float64
  HighPrice float64
  LowPrice float64
  PriceChangePercent1h float64
  PriceChangePercent24h float64
  PriceChangePercent7d float64
  Volume float64
  Spread float64
}

// Balance
type Balance struct {
  ServiceName string
  CurrencyCode string
  Amount float64
  Available float64
}
type BalanceArray []Balance

func (barr BalanceArray) Extract(currency_code string) (b Balance){
  for _, b = range barr {
    if b.CurrencyCode == currency_code {
      return
    }
  }
  return
}

func (barr1 BalanceArray) DropZero() (barr2 BalanceArray) {
  for _, b := range barr1 {
    if b.Amount > 0 {
      barr2 = append(barr2, b)
    }
  }
  return
}

//Order info
type OrderInfo struct{
  OrderId string
  CurrencyPair string
  Side string
  Type string
  Status string
  TotalSize float64
  ExecutedSize float64
  Price float64
  AveragePrice float64
}

//AvailablePairs
type AvailablePairs []string
func (ms AvailablePairs) Have (currency_pair string, service string) bool{
  switch service{
  case "bitflyer":
    currency_pair = strings.ToUpper(currency_pair)
  case "bitbank":
    currency_pair = strings.ToLower(currency_pair)
  case "binance":
    currency_pair = strings.ToUpper(strings.Replace(currency_pair, "_", "", -1))
  }
  for _, m := range ms {
    if m == currency_pair {
      return true
    }
  }
  return false
}

// for calculate
func Round(f float64, places int) (float64) {
    shift := math.Pow(10, float64(places))
    return math.Floor(f * shift + .5) / shift
}
