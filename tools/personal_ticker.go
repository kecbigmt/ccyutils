package tools

import (
  "fmt"
  "log"
  "errors"

  bn "github.com/kecbigmt/ccyutils/binance"
  bf "github.com/kecbigmt/ccyutils/bitflyer"
  "github.com/kecbigmt/ccyutils"
)

type PersonalTick struct{
  ServiceName string
  CurrencyCode string
  KeyCurrencyCode string
  Amount float32
  Amount_BTC float32
  Amount_key float32
  LastPrice_BTC float32
  LastPrice_key float32
  Spread float32
  Volume float32
}
type PersonalTickArray []PersonalTick
func (p PersonalTickArray) SumBTC() (sum float32) {
  for _, t := range p{
    sum += t.Amount_BTC
  }
  return
}
func (p PersonalTickArray) SumKey() (sum float32) {
  for _, t := range p{
    sum += t.Amount_key
  }
  return
}

// 手持ちの通貨の残高とTicker情報を取得
func PersonalTicker(key_currency_code string) (ptarr PersonalTickArray){
  // まずBTCJPYを取得
  key_tick, err := bf.Ticker("BTC_"+key_currency_code)
  if err != nil{
    log.Fatal(err)
  }
  key_ltp := key_tick.LastPrice
  bf_balances, _ := bf.Balance()
  bf_balances = bf_balances.DropZero()
  bn_balances, _ := bn.Balance()
  bn_balances = bn_balances.DropZero()
  balances := append(bf_balances, bn_balances...)
  for _, b := range balances{
    ptarr = append(ptarr, BalanceToTicker(b, key_ltp, key_currency_code))
  }
  return
}

func BalanceToTicker(b ccyutils.Balance, key_ltp float32, key_currency_code string) (pt PersonalTick){
  currency_pair := b.CurrencyCode+"_BTC"
  // キー通貨もしくはBTCであればAPIを叩く必要がない
  if b.CurrencyCode == key_currency_code {
    pt.ServiceName = b.ServiceName
    pt.CurrencyCode = b.CurrencyCode
    pt.Amount = b.Amount
    pt.LastPrice_BTC = 1.0/key_ltp
    pt.Amount_BTC = b.Amount*pt.LastPrice_BTC
    pt.LastPrice_key = 1.0
    pt.Amount_key = b.Amount
    pt.KeyCurrencyCode = key_currency_code
    return
  }
  if b.CurrencyCode == "BTC" {
    pt.ServiceName = b.ServiceName
    pt.CurrencyCode = b.CurrencyCode
    pt.Amount = b.Amount
    pt.LastPrice_BTC = 1.0/key_ltp
    pt.Amount_BTC = b.Amount
    pt.LastPrice_key = key_ltp
    pt.Amount_key = b.Amount * key_ltp
    pt.KeyCurrencyCode = key_currency_code
    return
  }

  // キー通貨でもBTCでもなかったらTicker情報を取得
  var tick ccyutils.Tick
  switch b.ServiceName{
  case "bitflyer":
    tick, _ = bf.Ticker(currency_pair)
  case "binance":
    tick, _ = bn.Ticker(currency_pair)
  default:
    err := errors.New(fmt.Sprintf("[Error]%v not found", b.ServiceName))
    log.Fatal(err)
  }
  pt.ServiceName = b.ServiceName
  pt.CurrencyCode = b.CurrencyCode
  pt.LastPrice_BTC = tick.LastPrice
  pt.LastPrice_key = pt.LastPrice_BTC * key_ltp
  pt.Amount = b.Amount
  pt.Amount_BTC = b.Amount*pt.LastPrice_BTC
  pt.Amount_key = b.Amount*pt.LastPrice_key
  pt.KeyCurrencyCode = key_currency_code
  pt.Spread = tick.Spread
  pt.Volume = tick.Volume
  return
}
