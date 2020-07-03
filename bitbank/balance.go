package bitbank

import (
  "fmt"
  "strings"
  "strconv"
  "io/ioutil"
  "encoding/json"

  "github.com/kecbigmt/ccyutils"
)

type BitBankBalanceResp struct {
  Success int `json:"success"`
  Data struct {
    Assets BitBankBalanceArray `json:"assets"`
  } `json:"data"`
}
type BitBankBalance struct {
  Asset string `json:"asset"`
  AmountPrecision int `json:"amount_precision"`
  OnhandAmount string `json:"onhand_amount"`
  LockedAmount string `json:"locked_amount"`
  FreeAmount string `json:"free_amount"`
  StopDeposit bool `json:"stop_deposit"`
  StopWithdrawal bool `json:"stop_withdrawal"`
}
type BitBankBalanceArray []BitBankBalance
func (bbbarr BitBankBalanceArray) Norm() (barr ccyutils.BalanceArray) {
  for _, bbb := range bbbarr {
    onhand, _ := strconv.ParseFloat(bbb.OnhandAmount, 64)
    free, _ := strconv.ParseFloat(bbb.FreeAmount, 64)
    b := ccyutils.Balance{
      ServiceName: "bitbank",
      CurrencyCode: strings.ToUpper(bbb.Asset),
      Amount: onhand,
      Available: free,
    }
    barr = append(barr, b)
  }
  return
}

func Balance() (barr ccyutils.BalanceArray, err error) {
  resp, err := AuthorizedRequest("GET", "/user/assets", nil)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"300",
      "component":"balance",
      "service":"bitbank",
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
      "service":"bitbank",
      "message":"[Error]HTTP Error(%v)",
      "detail":%v
    }`, resp.StatusCode, string(bytes))
    return
  }

  var balance BitBankBalanceResp
  if err = json.Unmarshal(bytes, &balance); err != nil {
    err = fmt.Errorf(`{
      "error_code":"302",
      "component":"balance",
      "service":"bitbank",
      "message":"[Error]Failed to decode JSON",
      "detail":%v
    }`, err)
    return
  }
  switch balance.Success {
  case 1:
    barr = balance.Data.Assets.Norm()
    return
  default:
    err = fmt.Errorf(`{
      "error_code":"303",
      "component":"balance",
      "service":"bitbank",
      "message":"[Error]API Error.",
      "detail":%v
    }`, string(bytes))
    return
  }
}
