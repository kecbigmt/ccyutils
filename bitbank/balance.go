package bitbank

import (
  "fmt"
  "log"
  "errors"
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
  resp, err := AuthorizedRequest("GET", "/user/assets", map[string]string{})
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

  var balance BitBankBalanceResp
  if err := json.Unmarshal(bytes, &balance); err != nil {
      log.Fatal(err)
  }
  barr = balance.Data.Assets.Norm()
  return
}
