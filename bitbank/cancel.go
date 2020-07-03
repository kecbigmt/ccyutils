package bitbank

import (
  "fmt"
  "strings"
  "io/ioutil"
  "encoding/json"
)
type BitBankCancel struct{
  Success int `json:"success"`
}

func Cancel(currency_pair string, oid string) (bl bool, err error){
  t_currency_pair := strings.ToLower(currency_pair) // XXX_YYY -> xxx_yyy
  params := map[string]string{
    "pair": t_currency_pair,
    "order_id": oid,
  }
  resp, err := AuthorizedRequest("POST", "/user/spot/cancel_order", params)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"600",
      "component":"cancel",
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
      "error_code":"601",
      "component":"cancel",
      "service":"bitbank",
      "message":"[Error]HTTP Error(%v)",
      "detail":%v
    }`, resp.StatusCode, string(bytes))
    return
  }
  var bbc BitBankCancel
  if err = json.Unmarshal(bytes, &bbc); err != nil {
    err = fmt.Errorf(`{
      "error_code":"602",
      "component":"cancel",
      "service":"bitbank",
      "message":"[Error]Failed to decode JSON.",
      "detail":%v
    }`, err)
    return
  }
  switch bbc.Success{
  case 1:
    bl = true
    return
  default:
    bl = false
    err = fmt.Errorf(`{
      "error_code":"603",
      "component":"cancel",
      "service":"bitbank",
      "message":"[Error]API Error",
      "detail":%v
    }`, string(bytes))
    return
  }
}
