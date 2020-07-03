package tools

import (
  "fmt"
  "time"
  "strings"
  "net/url"
  "net/http"

  bn "github.com/kecbigmt/ccyutils/binance"
  //bf "github.com/kecbigmt/ccyutils/bitflyer"
  //bb "github.com/kecbigmt/ccyutils/bitbank"
  //"github.com/kecbigmt/ccyutils"
)

func Escape(service string, currency_pair string, size float64){
  switch service{
  case "bitflyer":
    return
  case "bitbank":
    return
  case "binance":
    for {
      oid, err := bn.SellMarket(currency_pair, size)
      if err != nil{
        fmt.Println(err)
      }else{
        fmt.Println(oid)
        line_notify(oid, "Nc2WtVS6AEnkpPdTmvd4TrnPhbmQUhXkzJRWOMcB1QT")
        break
      }
      time.Sleep(500*time.Millisecond)
    }
  }
}

func line_notify(text string, lineNotifyAccessToken string) {
  url_api := "https://notify-api.line.me/api/notify"
  client := &http.Client{Timeout: time.Duration(10 * time.Second)}

  strs := strings.Split(text, "\v")
  for _, t := range strs {
    values := url.Values{}
    values.Set("message", t)
    req, err := http.NewRequest("POST", url_api, strings.NewReader(values.Encode()))
    if err != nil {
      fmt.Println(err)
      return
    }
    req.Header.Set("Authorization", "Bearer " + lineNotifyAccessToken)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, err := client.Do(req)
    if err != nil {
      fmt.Println(err)
      return
    }
    fmt.Println(resp.Status)
    resp.Body.Close()
  }
}
