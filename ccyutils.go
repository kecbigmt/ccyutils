// common
package ccyutils
import "math"

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

// for calculate
func Round(f float64, places int) (float64) {
    shift := math.Pow(10, float64(places))
    return math.Floor(f * shift + .5) / shift
}
