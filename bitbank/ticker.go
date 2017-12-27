package bitbank

import (
  "log"
  "errors"
  "strings"
  "strconv"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "github.com/kecbigmt/ccyutils"
)

type BitBankTick struct {
  Success int `json:"success"`
  Data BitBankData `json:"data"`
}
type BitBankData struct {
  Sell string `json:"sell"`
  Buy string `json:"buy"`
  High string `json:"high"`
  Low string `json:"Low"`
  Last string `json:"last"`
  Vol string `json:"vol"`
  Timestamp int64 `json:"timestamp"`
}

func (bbt BitBankTick) Norm(currency_pair string) ccyutils.Tick{
  var td ccyutils.Tick
  td.ServiceName = "bitbank"
  td.CurrencyPair = currency_pair
  td.UnixTimestamp = bbt.Data.Timestamp / int64(1000)
  b, _ := strconv.ParseFloat(bbt.Data.Buy, 64)
  td.BestBid = b
  s, _ := strconv.ParseFloat(bbt.Data.Sell, 64)
  td.BestAsk = s
  lt, _ := strconv.ParseFloat(bbt.Data.Last, 64)
  td.LastPrice = lt
  h, _ := strconv.ParseFloat(bbt.Data.High, 64)
  td.HighPrice = h
  lw, _ := strconv.ParseFloat(bbt.Data.Low, 64)
  td.LowPrice = lw
  v, _ := strconv.ParseFloat(bbt.Data.Vol, 64)
  td.Volume = v
  td.Spread = (td.BestAsk - td.BestBid) / td.BestAsk
  return td
}

func Ticker(currency_pair string) (tick ccyutils.Tick, err error){
  t_currency_pair := strings.ToLower(currency_pair) // XXX_YYY -> xxx_yyy
  url := "https://public.bitbank.cc/"+t_currency_pair+"/ticker"
  req, _ := http.NewRequest("GET", url, nil)
  client := new(http.Client)
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil{
    return
  }

  bytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      log.Fatal(err)
  }
  var bbt BitBankTick
  if err := json.Unmarshal(bytes, &bbt); err != nil {
      log.Fatal(err)
  }
  if bbt.Success == 0{
    err = errors.New("[Error]API Error")
    return
  }
  tick = bbt.Norm(currency_pair)
  return
}
