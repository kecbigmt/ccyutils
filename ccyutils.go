// common
package ccyutils
import "math"

// Ticker
type Tick struct {
  ServiceName string
  CurrencyPair string
  UnixTimestamp int64
  TickId int
  BestBid float32
  BestAsk float32
  BestBidSize float32
  BestAskSize float32
  TotalBidDepth float32
  TotalAskDepth float32
  LastPrice float32
  HighPrice float32
  LowPrice float32
  PriceChangePercent1h float32
  PriceChangePercent24h float32
  PriceChangePercent7d float32
  Volume float32
  Spread float32
}

// Balance
type Balance struct {
  ServiceName string
  CurrencyCode string
  Amount float32
  Available float32
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
