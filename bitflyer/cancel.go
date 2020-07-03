package bitflyer

import (
  "fmt"
  "strings"
)

func Cancel(currency_pair string, oid string) (bl bool, err error){
  t_currency_pair := strings.ToUpper(currency_pair) // XXX_YYY
  params := map[string]string{
    "product_code": t_currency_pair,
    "child_order_acceptance_id": oid,
  }
  resp, err := AuthorizedRequest("POST", "/v1/me/cancelchildorder", params)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"600",
      "component":"cancel",
      "service":"bitflyer",
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
      "service":"bitflyer",
      "message":"[Error]HTTP Error(%v)"
      "detail":""
    }`, resp.StatusCode)
    return
  }
}
