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

type PersonalTick struct{
  ServiceName string
  CurrencyCode string
  KeyCurrencyCode string
  Amount float64
  Amount_BTC float64
  Amount_key float64
  LastPrice_BTC float64
  LastPrice_key float64
  Spread float64
  Volume float64
}
type PersonalTickArray []PersonalTick
func (p PersonalTickArray) SumBTC() (sum float64) {
  for _, t := range p{
    sum += t.Amount_BTC
  }
  return
}
func (p PersonalTickArray) SumKey() (sum float64) {
  for _, t := range p{
    sum += t.Amount_key
  }
  return
}
func (ptarr PersonalTickArray) Search(currency_code string, service_i interface{}) (output PersonalTick, ok bool){
  switch service := service_i.(type){
  case string:
    for _, pt := range ptarr {
      if pt.ServiceName == service && pt.CurrencyCode == currency_code{
        output = pt
        ok = true
        return
      }
    }
    ok = false
    return
  default:
    for _, pt := range ptarr {
      if pt.CurrencyCode == currency_code{
        output = pt
        ok = true
        return
      }
    }
    ok = false
    return
  }
}
//
type AvailablePairsStruct struct{
  BitFlyer ccyutils.AvailablePairs
  BitBank ccyutils.AvailablePairs
  Binance ccyutils.AvailablePairs
}
func GetAllAvailablePairs() AvailablePairsStruct{
  return AvailablePairsStruct{
    BitFlyer: bf.AvailablePair(),
    BitBank: bb.AvailablePair(),
    Binance: bn.AvailablePair(),
  }
}



// 手持ちの通貨の残高とTicker情報を取得
func PersonalTicker(key_currency_code string, pairs AvailablePairsStruct) (arr PersonalTickArray){
  // 通貨ペアを取得
  bitflyer_pairs := pairs.BitFlyer
  bitbank_pairs := pairs.BitBank
  binance_pairs := pairs.Binance
  // まずBTCJPYを取得（bitFlyerの相場を基準にする）
  /*key_tick, err := bf.Ticker("BTC_"+key_currency_code)
  if err != nil{
    log.Fatal(err)
  }
  key_ltp := key_tick.LastPrice*/
  arr = append(arr, BalanceToTickerMulti(bf.Balance, key_currency_code, &bitflyer_pairs)...)
  arr = append(arr, BalanceToTickerMulti(bb.Balance, key_currency_code, &bitbank_pairs)...)
  arr = append(arr, BalanceToTickerMulti(bn.Balance, key_currency_code, &binance_pairs)...)
  // 0.0になっている値があれば計算用のLTPを決定し、なければ処理終了
  var ltp_btckey float64
  for _, a := range arr {
    if a.LastPrice_BTC == 0.0 || a.LastPrice_key == 0.0{
      // BTC/JPYを取得済みデータのなかから探し、なければbitFlyerのAPIを叩く
      if tick_btc, ok := arr.Search("BTC", nil); ok{
        ltp_btckey = tick_btc.LastPrice_key
      } else {
        key_tick, err := bf.Ticker("BTC_"+key_currency_code)
        if err != nil{
          log.Fatal(err)
        }
        ltp_btckey = key_tick.LastPrice
      }
      break
    }
  }
  if ltp_btckey == 0.0 {
    return
  }

  // 計算用のLTPを用いて0.0を埋める
  for i, a := range arr {
    switch {
    case a.CurrencyCode == "BTC" && a.LastPrice_key == 0.0:
      (arr[i]).LastPrice_BTC = 1.0
      (arr[i]).LastPrice_key = ltp_btckey
      (arr[i]).Amount_BTC = (arr[i]).Amount
      (arr[i]).Amount_key = (arr[i]).Amount * (arr[i]).LastPrice_key
    case a.LastPrice_BTC == 0.0:
      (arr[i]).LastPrice_BTC = (arr[i]).LastPrice_key * (1.0/ltp_btckey)
      (arr[i]).Amount_BTC = (arr[i]).Amount * (arr[i]).LastPrice_BTC
    case a.LastPrice_key == 0.0:
      (arr[i]).LastPrice_key = (arr[i]).LastPrice_BTC * ltp_btckey
      (arr[i]).Amount_key = (arr[i]).Amount * (arr[i]).LastPrice_key
    default:
      break
    }
  }
  return
}

func BalanceToTickerMulti (f func()(ccyutils.BalanceArray, error), key_currency_code string, available_pairs *ccyutils.AvailablePairs) (arr PersonalTickArray){
  balances, _ := f()
  balances = balances.DropZero()
  for _, b := range balances{
    arr = append(arr, BalanceToTicker(&b, key_currency_code, available_pairs))
  }
  // 0.0になっているLTPを同じ取引所のBTC/JPYで計算する
  for i, a := range arr {
    if a.LastPrice_BTC == 0.0 && a.LastPrice_key == 0.0{
      return
    }
    if a.LastPrice_BTC == 0.0 || a.LastPrice_key == 0.0{
      // BTC/JPYを取得済みデータのなかから探し、なければ後回し
      if tick_btc, ok := arr.Search("BTC", a.ServiceName); !ok{
        return
      }else{
        ltp_btckey := tick_btc.LastPrice_key
        switch {
        case a.LastPrice_BTC == 0.0 :
          (arr[i]).LastPrice_BTC = (arr[i]).LastPrice_key * (1.0/ltp_btckey)
          (arr[i]).Amount_BTC = (arr[i]).Amount * (arr[i]).LastPrice_BTC
        case a.LastPrice_key == 0.0 :
          (arr[i]).LastPrice_key = (arr[i]).LastPrice_BTC * ltp_btckey
          (arr[i]).Amount_key = (arr[i]).Amount * (arr[i]).LastPrice_key
        }
      }
    }
  }
  return
}

func BalanceToTicker(b *ccyutils.Balance, key_currency_code string, available_pairs *ccyutils.AvailablePairs) (pt PersonalTick){
  // decide currency_pair
  var currency_pair string
  var currency_opposite string
  switch{
  case available_pairs.Have(b.CurrencyCode+"_BTC", b.ServiceName):
    currency_opposite = "_BTC"
    currency_pair = b.CurrencyCode+currency_opposite
  case available_pairs.Have("BTC_"+b.CurrencyCode, b.ServiceName):
    currency_opposite = "BTC_"
    currency_pair = currency_opposite+b.CurrencyCode
  case available_pairs.Have(b.CurrencyCode+"_"+key_currency_code, b.ServiceName):
    currency_opposite = "_"+key_currency_code
    currency_pair = b.CurrencyCode+currency_opposite
  default:
    pt.ServiceName = b.ServiceName
    pt.CurrencyCode = b.CurrencyCode
    pt.Amount = b.Amount
    pt.KeyCurrencyCode = key_currency_code
    return
  }
  // get ticker
  var tick ccyutils.Tick
  var err error
  switch b.ServiceName{
  case "bitflyer":
    tick, _ = bf.Ticker(currency_pair)
  case "binance":
    tick, _ = bn.Ticker(currency_pair)
  case "bitbank":
    tick, _ = bb.Ticker(currency_pair)
  default:
    err = errors.New(fmt.Sprintf("[Error]%v not found", b.ServiceName))
    log.Fatal(err)
  }
  // data
  switch currency_opposite{
  case "_BTC":
    pt.LastPrice_BTC = tick.LastPrice
  case "BTC_":
    pt.LastPrice_BTC = 1.0/tick.LastPrice
  case "_"+key_currency_code:
    pt.LastPrice_key = tick.LastPrice
  }
  switch b.CurrencyCode{
  case "BTC":
    pt.LastPrice_BTC = 1.0
    pt.Spread = tick.Spread
    pt.Volume = tick.Volume
  case key_currency_code:
    pt.LastPrice_key = 1.0
    pt.Spread = 0.0
    pt.Volume = 0.0
  default:
    pt.Spread = tick.Spread
    pt.Volume = tick.Volume
  }
  pt.ServiceName = b.ServiceName
  pt.CurrencyCode = b.CurrencyCode
  pt.Amount = b.Amount
  pt.Amount_BTC = b.Amount*pt.LastPrice_BTC
  pt.Amount_key = b.Amount*pt.LastPrice_key
  pt.KeyCurrencyCode = key_currency_code
  return
}
