package bitflyer

import (
  "fmt"
  "io/ioutil"
  "encoding/json"

  "github.com/kecbigmt/ccyutils"
)

type BitFlyerBalance struct {
    CurrencyCode string `json:"currency_code"`
    Amount float64 `json:"amount"`
    Available float64 `json:"available"`
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
  resp, err := AuthorizedRequest("GET", "/v1/me/getbalance", nil)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"300",
      "component":"balance",
      "service":"bitflyer",
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
      "service":"bitflyer",
      "message":"[Error]HTTP Error(%v)",
      "detail":%v
    }`, resp.StatusCode, string(bytes))
    return
  }

  var balances BitFlyerBalanceArray
  if err = json.Unmarshal(bytes, &balances); err != nil {
    err = fmt.Errorf(`{
      "error_code":"302",
      "component":"balance",
      "service":"bitflyer",
      "message":"[Error]Failed to decode JSON",
      "detail":%v
    }`, err)
    return
  }
  barr = balances.Norm()
  return
}
