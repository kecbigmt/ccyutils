// common
package ccyutils

type TickData struct {
  ExchangeName string
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
}
