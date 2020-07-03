package binance

import (
  "fmt"
  "strconv"
  "io/ioutil"
  "encoding/json"

  "github.com/kecbigmt/ccyutils"
)

type BinanceAccount struct {
  MakerCommission int `json:"makerCommission"`
  TakerCommission int `json:"takerCommission"`
  BuyerCommission int `json:"buyerCommission"`
  SellerCommission int `json:"sellerCommission"`
  CanTrade bool `json:"canTrade"`
  CanWithdraw bool `json:"canWithdraw"`
  CanDeposit bool `json:"canDeposit"`
  UpdateTime int64 `json:"updateTime"`
  Balances BinanceBalanceArray `json:"balances"`
}
type BinanceBalance struct {
  Asset string `json:"asset"`
  Free string `json:"free"`
  Locked string `json:"locked"`
}
type BinanceBalanceArray []BinanceBalance
func (bbarr BinanceBalanceArray) Norm() (barr ccyutils.BalanceArray) {
  for _, bb := range bbarr {
    free, _ := strconv.ParseFloat(bb.Free, 64)
    locked, _ := strconv.ParseFloat(bb.Locked, 64)
    b := ccyutils.Balance{
      ServiceName: "binance",
      CurrencyCode: bb.Asset,
      Amount: free+locked,
      Available: free,
    }
    barr = append(barr, b)
  }
  return
}

func Balance() (barr ccyutils.BalanceArray, err error) {
  resp, err := AuthorizedRequest("GET", "/api/v3/account", nil, true)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"300",
      "component":"balance",
      "service":"binance",
      "message":"[Error]Failed to request",
      "detail":%v
      }`, err)
    return
  }
  defer resp.Body.Close()
  bytes, _ := ioutil.ReadAll(resp.Body)
  if resp.StatusCode != 200{
    err = fmt.Errorf(`{
      "error_code":"301",
      "component":"balance",
      "service":"binance",
      "message":"[Error]HTTP Error(%v)",
      "detail":%v
      }`, resp.StatusCode, string(bytes))
    return
  }

  var account BinanceAccount
  if err = json.Unmarshal(bytes, &account); err != nil {
    err = fmt.Errorf(`{
      "error_code":"302",
      "component":"balance",
      "service":"binance",
      "message":"[Error]Failed to decode JSON",
      "detail":%v
      }`, err)
    return
  }
  barr = account.Balances.Norm()
  return
}
