package tools

import(
  "fmt"
  "testing"
  //bn "github.com/kecbigmt/ccyutils/binance"
  bf "github.com/kecbigmt/ccyutils/bitflyer"
  //bb "github.com/kecbigmt/ccyutils/bitbank"
)

/*func TestAvailablePairs(t *testing.T){
  fmt.Println(PersonalTicker("JPY"))
}*/

func BenchmarkNormal(b *testing.B){
  b.ReportAllocs()
  key_currency_code := "JPY"
  // 通貨ペアを取得
  available_pairs, _ := bf.AvailablePair()
  // 通貨ペア辞書の準備
  pair_dict := make(map[string]string)

  balances, _ := bf.Balance()
  balances = balances.DropZero()
  for i := 0; i<=100; i++{
    for _, b := range balances{
      var currency_pair string
      currency_opposite, ok := (pair_dict)[b.CurrencyCode]
      if !ok {
        switch{
        case available_pairs.Have(b.CurrencyCode+"_BTC", b.ServiceName):
          (pair_dict)[b.CurrencyCode] = "BTC"
          currency_opposite = "BTC"
        case available_pairs.Have(b.CurrencyCode+"_"+key_currency_code, b.ServiceName):
          (pair_dict)[b.CurrencyCode] = key_currency_code
          currency_opposite = key_currency_code
        default:
          return
        }
      }
      currency_pair = b.CurrencyCode+"_"+currency_opposite
      fmt.Sprintf("%v", currency_pair)
      return
    }
  }
}

func BenchmarkSimple(b *testing.B){
  b.ReportAllocs()
  key_currency_code := "JPY"
  // 通貨ペアを取得
  available_pairs, _ := bf.AvailablePair()

  balances, _ := bf.Balance()
  balances = balances.DropZero()
  for i := 0; i<=100; i++{
    for _, b := range balances{
      var currency_pair string
      var currency_opposite string
      switch{
      case available_pairs.Have(b.CurrencyCode+"_BTC", b.ServiceName):
        currency_opposite = "BTC"
      case available_pairs.Have(b.CurrencyCode+"_"+key_currency_code, b.ServiceName):
        currency_opposite = key_currency_code
      default:
        return
      }
      currency_pair = b.CurrencyCode+"_"+currency_opposite
      fmt.Sprintf("%v", currency_pair)
      return
    }
  }
}
