package binance

import (
  "fmt"
  "strconv"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type BinanceOrdered struct {
  Symbol string `json:"symbol"`
  OrderId int `json:"orderId"`
  ClientOrderId string `json:"clientOrderId"`
  TransactTime int `json:"transactTime"`
}

func Order(currency_pair string, order_type string, side string, size float64, price float64) (oid string, err error){
  t_currency_pair := strings.ToUpper(strings.Replace(currency_pair, "_", "", -1)) // XXX_YYY -> XXXYYY
  order_type = strings.ToUpper(order_type)
  side = strings.ToUpper(side)
  var params map[string]string
  switch order_type{
  case "LIMIT":
    params = map[string]string{
      "symbol": t_currency_pair,
      "type": order_type,
      "side": side,
      "price": strconv.FormatFloat(price, 'f', 8, 64),
      "quantity": strconv.FormatFloat(size, 'f', 8, 64),
      "timeInForce": "GTC",
    }
  case "MARKET":
    params = map[string]string{
      "symbol": t_currency_pair,
      "type": order_type,
      "side": side,
      "quantity": strconv.FormatFloat(size, 'f', 8, 64),
    }
  }
  resp, err := AuthorizedRequest("POST", "/api/v3/order", params, true)
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"500",
      "component":"order",
      "service":"binance",
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
      "service":"binance",
      "message":"[Error]HTTP Error(%v)",
      "detail":%v
    }`, resp.StatusCode, string(bytes))
    return
  }
  var bo BinanceOrdered
  if err = json.Unmarshal(bytes, &bo); err != nil {
    err = fmt.Errorf(`{
      "error_code":"502",
      "component":"order",
      "service":"binance",
      "message":"[Error]Failed to decode JSON",
      "detail":%v
    }`, err)
    return
  }
  oid = strconv.Itoa(bo.OrderId)
  return
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
