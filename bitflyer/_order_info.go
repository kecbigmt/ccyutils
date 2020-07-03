package bitflyer

import (
  "fmt"
  "strconv"
  "strings"
  "io/ioutil"
  "encoding/json"

  "github.com/kecbigmt/ccyutils"
)

type BitFlyerOrderInfo struct{
  ChildOrderId string `json:"child_order_id"`
  ProductCode string `json:"product_code"`
  ChildOrderState string `json:"child_order_state"`
  OutstandingSize int `json:"outstanding_size"`
  CancelSize int `json:"cancel_size"`
  ExecutedSize int `json:"executed_size"`
  TotalCommission int `json:"total_commission"`
  Side string `json:"side"`
  ChildOrderType `json:"child_order_type"`
  Price float64 `json:"price"`
  AveragePrice float64 `json:"average_price"`
}

type BitFlyerOrders []OrderInfo
func (os BitFlyerOrders) Norm(currency_pair string, order_id string) (info ccyutils.OrderInfo, err error){
  for _, bfoi := range bfo{
    if bfoi.ChildOrderId == oid {
      info.OrderId = oid
      info.CurrencyPair = currency_pair
      info.Side = bfoi.Side
      info.Type = bfoi.ChildOrderType
      info.Status = bfoi.ChildOrderState
      info.TotalSize = bfoi.TotalCommission
      info.ExecutedSize = bfoi.ExecutedSize
      info.Price = bfoi.Price
      info.AveragePrice = bfoi.AveragePrice
      return
    }
  }
  err = fmt.Errorf(`{
    "error_code":"700",
    "component":"order_info",
    "service":"bitflyer",
    "message":"[Error]Order ID '%v' not found.",
    "detail":""
  }`, oid)
  return
}

func OrderInfo(currency_pair string, order_id string) (info ccyutils.OrderInfo, err error){
  t_currency_pair := strings.ToUpper(currency_pair) // XXX_YYY
  params := map[string]string{
    "product_code": t_currency_pair,
    "child_order_id": order_id,
  }
  resp, _ := AuthorizedRequest("POST", "/v1/me/sendchildorder", params)
  defer resp.Body.Close()
  if err != nil{
    err = fmt.Errorf(`{
      "error_code":"700",
      "component":"order_info",
      "service":"bitflyer",
      "message":"[Error]Failed to request.",
      "detail":%v
    }`, err)
    return
  }
  bytes, _ := ioutil.ReadAll(resp.Body)
  if resp.StatusCode != 200{
    err = fmt.Errorf(`{
      "error_code":"701",
      "component":"order_info",
      "service":"bitflyer",
      "message":"[Error]HTTP Error(%v)",
      "detail":%v
    }`, resp.StatusCode, string(bytes))
    return
  }
  var bfo BitFlyerOrders
  if err = json.Unmarshal(bytes, &bfo); err != nil {
    err = fmt.Errorf(`{
      "error_code":"702",
      "component":"order_info",
      "service":"bitflyer",
      "message":"[Error]Failed to decode JSON",
      "detail":%v
    }`, err)
    return
  }
  info, err = bfo.Norm(currency_pair string, order_id string)
  return
}
