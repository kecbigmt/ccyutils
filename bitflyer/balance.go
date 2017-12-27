package bitflyer

import (
  "fmt"
  "log"
  "io/ioutil"
  "encoding/json"

  "github.com/kecbigmt/ccyutils"
)

type BitFlyerBalance struct {
    CurrencyCode string `json:"currency_code"`
    Amount float32 `json:"amount"`
    Available float32 `json:"available"`
}
type BitFlyerBalanceArray []BitFlyerBalance
func (bfbarr BitFlyerBalanceArray) Norm() (barr ccyutils.BalanceArray) {
  for _, bfb := range bfbarr {
    b := ccyutils.Balance{
      ServiceName: "bitflyer",
      CurrencyCode: bfb.CurrencyCode,
      Amount: bfb.Amount,
      Available: bfb.Available,
    }
    barr = append(barr, b)
  }
  return
}

func Balance() (barr ccyutils.BalanceArray, err error) {
  resp, err := AuthorizedRequest("GET", "/v1/me/getbalance", map[string]string{})
  defer resp.Body.Close()
  if err != nil {
      return
  }

  bytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      log.Fatal(err)
  }
  if resp.StatusCode != 200 {
    err := fmt.Sprintf("HTTP Error(%v): %v", resp.StatusCode, string(bytes))
    log.Fatal(err)
  }

  var balances BitFlyerBalanceArray
  if err := json.Unmarshal(bytes, &balances); err != nil {
      log.Fatal(err)
  }
  barr = balances.Norm()
  return
}
