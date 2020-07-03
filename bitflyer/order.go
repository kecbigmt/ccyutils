package bitflyer

import (
  "fmt"
  "strconv"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type BitFlyerOrdered struct {
  OrderId string `json:"child_order_acceptance_id"`
}

func Order(currency_pair string, child_order_type string, side string, size float64, price float64, minute_to_expire int, time_in_force string) (oid string, err error){
  t_currency_pair := strings.ToUpper(currency_pair) // XXX_YYY
  params := map[string]string{
    "product_code": t_currency_pair,
    "child_order_type": child_order_type,
    "side": side,
    "price": strconv.FormatFloat(price, 'f', 8, 64),
    "size": strconv.FormatFloat(size, 'f', 8, 64),
    "minute_to_expire": string(minute_to_expire),
    "time_in_force": time_in_force,
  }
  resp, _ := AuthorizedRequest("POST", "/v1/me/sendchildorder", params)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"500",
      "component":"order",
      "service":"bitflyer",
      "message":"[Error]Failed to request)",
      "detail":%v
    }`, err)
    return
  }
  defer resp.Body.Close()
  bytes, _ := ioutil.ReadAll(resp.Body)
  if resp.StatusCode != 200{
    err = fmt.Errorf(`{
      "error_code":"501",
      "component":"order",
      "service":"bitflyer",
      "message":"[Error]HTTP Error(%v)",
      "detail":%v
    }`, resp.StatusCode, string(bytes))
    return
  }
  var bfo BitFlyerOrdered
  if err = json.Unmarshal(bytes, &bfo); err != nil {
    err = fmt.Errorf(`{
      "error_code":"502",
      "component":"order",
      "service":"bitflyer",
      "message":"[Error]Failed to decode JSON",
      "detail":%v
    }`, err)
    return
  }
  oid = bfo.OrderId
  return
}

func BuyLimit(currency_pair string, price float64, size float64) (oid string, err error){
  oid, err = Order(currency_pair, "LIMIT", "BUY", size, price, 43200, "GTC")
  return
}

func SellLimit(currency_pair string, price float64, size float64) (oid string, err error){
  oid, err = Order(currency_pair, "LIMIT", "SELL", size, price, 43200, "GTC")
  return
}

func BuyMarket(currency_pair string, size float64) (oid string, err error){
  var price float64
  oid, err = Order(currency_pair, "MARKET", "BUY", size, price, 43200, "GTC")
  return
}

func SellMarket(currency_pair string, size float64) (oid string, err error){
  var price float64
  oid, err = Order(currency_pair, "MARKET", "SELL", size, price, 43200, "GTC")
  return
}
