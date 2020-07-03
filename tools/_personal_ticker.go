package tools

import (
  "fmt"
  "log"
  "errors"

  bn "github.com/kecbigmt/ccyutils/binance"
  bf "github.com/kecbigmt/ccyutils/bitflyer"
  bb "github.com/kecbigmt/ccyutils/bitbank"
  "github.com/kecbigmt/ccyutils"
)

// 手持ちの通貨の残高とTicker情報を取得
func PersonalTicker(key_currency_code string) (ptarr PersonalTickArray){
  // まずBTCJPYを取得（bitFlyerの相場を基準にする）
  key_tick, err := bf.Ticker("BTC_"+key_currency_code)
  if err != nil{
    log.Fatal(err)
  }
  key_ltp := key_tick.LastPrice
  bf_balances, _ := bf.Balance()
  bf_balances = bf_balances.DropZero()
  bn_balances, _ := bn.Balance()
  bn_balances = bn_balances.DropZero()
  bb_balances, _ := bb.Balance()
  bb_balances = bb_balances.DropZero()
  balances := append(bf_balances, bn_balances...)
  balances = append(balances, bb_balances...)
  for _, b := range balances{
    ptarr = append(ptarr, BalanceToTicker(b, key_ltp, key_currency_code))
  }
  return
}

func BalanceToTicker(b ccyutils.Balance, key_ltp float64, key_currency_code string) (pt PersonalTick){
  currency_pair := b.CurrencyCode+"_BTC"
  // キー通貨もしくはBTCであればAPIを叩く必要がない
  if b.CurrencyCode == key_currency_code {
    pt.ServiceName = b.ServiceName
    pt.CurrencyCode = b.CurrencyCode
    pt.Amount = b.Amount
    pt.LastPrice_BTC = (1.0/key_ltp)
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
    pt.LastPrice_BTC = 1.0
    pt.Amount_BTC = b.Amount
    pt.LastPrice_key = key_ltp
    pt.Amount_key = b.Amount * key_ltp
    pt.KeyCurrencyCode = key_currency_code
    return
  }

  // キー通貨でもBTCでもなかったらTicker情報を取得
  var tick ccyutils.Tick
  var err error
  if b.CurrencyCode == "BTC"{

  }else{

    }
  switch b.ServiceName{
  case "bitflyer":
    tick, _ = bf.Ticker(currency_pair)
    pt.LastPrice_BTC = tick.LastPrice
    pt.LastPrice_key = pt.LastPrice_BTC * key_ltp
  case "binance":
    tick, _ = bn.Ticker(currency_pair)
    pt.LastPrice_BTC = tick.LastPrice
    pt.LastPrice_key = pt.LastPrice_BTC * key_ltp
  case "bitbank":
    tick, err = bb.Ticker(currency_pair)
    // 対BTCで取得できなければ対キー通貨（JPY）でも取得してみる
    if err != nil {
      tick, err = bb.Ticker(b.CurrencyCode+"_"+key_currency_code)
      if err != nil{
        log.Fatal(err)
      }else{
        pt.LastPrice_key = tick.LastPrice
        pt.LastPrice_BTC = pt.LastPrice_key * 1.0/key_ltp
      }
    } else {
      pt.LastPrice_BTC = tick.LastPrice
      pt.LastPrice_key = pt.LastPrice_key * key_ltp
    }
  default:
    err := errors.New(fmt.Sprintf("[Error]%v not found", b.ServiceName))
    log.Fatal(err)
  }
  pt.ServiceName = b.ServiceName
  pt.CurrencyCode = b.CurrencyCode
  pt.Amount = b.Amount
  pt.Amount_BTC = b.Amount*pt.LastPrice_BTC
  pt.Amount_key = b.Amount*pt.LastPrice_key
  pt.KeyCurrencyCode = key_currency_code
  pt.Spread = tick.Spread
  pt.Volume = tick.Volume
  return
}
