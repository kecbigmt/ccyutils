package bitbank

import (
    "log"
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

func (bbt BitBankTick) Norm(currency_pair string) ccyutils.TickData{
  var td ccyutils.TickData
  td.ExchangeName = "bitBank"
  td.CurrencyPair = currency_pair
  td.UnixTimestamp = bbt.Data.Timestamp / int64(1000)
  b, _ := strconv.ParseFloat(bbt.Data.Buy, 32)
  td.BestBid = float32(b)
  s, _ := strconv.ParseFloat(bbt.Data.Sell, 32)
  td.BestAsk = float32(s)
  lt, _ := strconv.ParseFloat(bbt.Data.Last, 32)
  td.LastPrice = float32(lt)
  h, _ := strconv.ParseFloat(bbt.Data.High, 32)
  td.HighPrice = float32(h)
  lw, _ := strconv.ParseFloat(bbt.Data.Low, 32)
  td.LowPrice = float32(lw)
  v, _ := strconv.ParseFloat(bbt.Data.Vol, 32)
  td.Volume = float32(v)
  return td
}

func Ticker(currency_pair string) ccyutils.TickData{
  t_currency_pair := strings.ToLower(currency_pair) // XXX_YYY -> xxx_yyy
  url := "https://public.bitbank.cc/"+t_currency_pair+"/ticker"
  req, _ := http.NewRequest("GET", url, nil)
  client := new(http.Client)
  resp, _ := client.Do(req)
  defer resp.Body.Close()

  bytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      log.Fatal(err)
  }
  var bbt BitBankTick
  if err := json.Unmarshal(bytes, &bbt); err != nil {
      log.Fatal(err)
  }
  return bbt.Norm(currency_pair)
}
