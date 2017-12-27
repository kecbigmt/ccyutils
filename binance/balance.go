package binance

import (
  "fmt"
  "log"
  "errors"
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
  resp, err := AuthorizedRequest("GET", "/api/v3/account", map[string]string{}, true)
  defer resp.Body.Close()
  if err != nil {
      return
  }

  bytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      log.Fatal(err)
  }
  if resp.StatusCode != 200 {
    err = errors.New(fmt.Sprintf("HTTP Error(%v): %v", resp.StatusCode, string(bytes)))
    return
  }

  var account BinanceAccount
  if err := json.Unmarshal(bytes, &account); err != nil {
      log.Fatal(err)
  }
  barr = account.Balances.Norm()
  return
}
