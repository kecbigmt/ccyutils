package bitflyer

import (
  "fmt"
  "log"
  "io/ioutil"
  "encoding/json"
)

type BalanceData struct {
    CurrencyCode string `json:"currency_code"`
    Amount float32 `json:"amount"`
    Available float32 `json:"available"`
}

func Balance() map[string]map[string]float32{
  resp := privateApi("GET", "/v1/me/getbalance", map[string]string{})
  defer resp.Body.Close()
  // JSONファイル読み込み
  bytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      log.Fatal(err)
  }
  if resp.StatusCode != 200 {
    err := fmt.Sprintf("HTTP Error(%v): %v", resp.StatusCode, string(bytes))
    log.Fatal(err)
  }

  // JSONデコード
  var balances []BalanceData
  if err := json.Unmarshal(bytes, &balances); err != nil {
      log.Fatal(err)
  }
  // デコードしたデータを表示
  var (
    jpy BalanceData
    btc BalanceData
    eth BalanceData
  )
  for _, j := range balances{
    switch j.CurrencyCode {
    case "JPY": jpy = j
    case "BTC": btc = j
    case "ETH": eth = j
    }
  }
  return map[string]map[string]float32{
    "JPY":map[string]float32{
      "amount": jpy.Amount,
      "available": jpy.Available,
    },
    "BTC":map[string]float32{
      "amount": btc.Amount,
      "available": btc.Available,
    },
    "ETH":map[string]float32{
      "amount": eth.Amount,
      "available": eth.Available,
    },
  }
}
