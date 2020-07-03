package binance

import (
  "fmt"
  "strings"
)

/*type BinanceCancel struct{
  Symbol string `json:"symbol"`
  OrigClientOrderId string `json:"origClientOrderId"`
  OrderId int `json:"orderId"`
  ClientOrderId string `json:"clientOrderId"`
}*/

func Cancel(currency_pair string, oid string) (bl bool, err error){
  t_currency_pair := strings.ToUpper(strings.Replace(currency_pair, "_", "", -1)) // XXX_YYY -> XXXYYY
  params := map[string]string{
    "symbol": t_currency_pair,
    "orderId": oid,
  }
  resp, err := AuthorizedRequest("DELETE", "/api/v3/order", params, true)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"600",
      "component":"cancel",
      "service":"binance",
      "message":"[Error]Failed to request.",
      "detail":%v
    }`, err)
    return
  }
  defer resp.Body.Close()
  switch resp.StatusCode{
  case 200:
    bl = true
    return
  default:
    bl = false
    err = fmt.Errorf(`{
      "error_code":"601",
      "component":"cancel",
      "service":"binance",
      "message":"[Error]HTTP Error(%v)"
      "detail":""
    }`, resp.StatusCode)
    return
  }
}
