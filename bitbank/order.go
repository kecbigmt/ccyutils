package bitbank

import (
  "fmt"
  "strconv"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type BitBankOrdered struct {
  Success int `json:"success"`
  Data struct {
    OrderId int `json:"order_id"`
    Pair string `json:"pair"`
    Side string `json:"side"`
    Type string `json:"type"`
    StartAmount string `json:"start_amount"`
    RemainingAmount string `json:"remaining_amount"`
    ExecutedAmount string `json:"executed_amount"`
    Price string `json:"price"`
    AveragePrice string `json:"average_price"`
    OrderedAt int `json:"ordered_at"`
    Status string `json:"status"`
  } `json:"data"`
  OrderId string `json:"child_order_acceptance_id"`
}

func Order(currency_pair string, order_type string, side string, size float64, price float64) (oid string, err error){
  t_currency_pair := strings.ToLower(currency_pair) // XXX_YYY -> xxx_yyy
  order_type = strings.ToLower(order_type)
  side = strings.ToLower(side)
  params := map[string]string{
    "pair": t_currency_pair,
    "type": order_type,
    "side": side,
    "price": strconv.FormatFloat(price, 'f', 8, 64),
    "amount": strconv.FormatFloat(size, 'f', 8, 64),
  }
  resp, err := AuthorizedRequest("POST", "/user/spot/order", params)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"500",
      "component":"order",
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
      "error_code":"501",
      "component":"order",
      "service":"bitbank",
      "message":"[Error]HTTP Error(%v)",
      "detail":%v
    }`, resp.StatusCode, string(bytes))
    return
  }
  var bbo BitBankOrdered
  if err = json.Unmarshal(bytes, &bbo); err != nil {
    err = fmt.Errorf(`{
      "error_code":"502",
      "component":"order",
      "service":"bitbank",
      "message":"[Error]Failed to decode JSON",
      "detail":%v
    }`, err)
    return
  }
  switch bbo.Success{
  case 1:
    oid = strconv.Itoa(bbo.Data.OrderId)
    return
  default:
    err = fmt.Errorf(`{
      "error_code":"503",
      "component":"order",
      "service":"bitbank",
      "message":"[Error]API Error.",
      "detail":%v
    }`, string(bytes))
    return
  }
}

func BuyLimit(currency_pair string, size float64, price float64) (oid string, err error){
  oid, err = Order(currency_pair, "LIMIT", "BUY", size, price)
  return
}

func SellLimit(currency_pair string, size float64, price float64) (oid string, err error){
  oid, err = Order(currency_pair, "LIMIT", "SELL", size, price)
  return
}

func BuyMarket(currency_pair string, size float64) (oid string, err error){
  var price float64
  oid, err = Order(currency_pair, "MARKET", "BUY", size, price)
  return
}

func SellMarket(currency_pair string, size float64) (oid string, err error){
  var price float64
  oid, err = Order(currency_pair, "MARKET", "SELL", size, price)
  return
}
