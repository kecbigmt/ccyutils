package binance

import (
  "os"
  "time"
  "strconv"
  "net/http"
  "net/url"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
)

func AuthorizedRequest(method string, path string, params map[string]string, sign bool) (resp *http.Response, err error){
  // init
  base_url := "https://api.binance.com"

  // prepare params
  var totalParams string
  values := url.Values{}
  if len(params)>0 {
    for key, value := range params{
      values.Add(key, value)
    }
  }
  if sign {
    timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(1000000), 10)
    values.Add("timestamp", timestamp)
    values.Add("recvWindow", "5000")
  }
  totalParams = values.Encode()

  // sign by hmac-sha256
  if sign {
    api_secret := os.Getenv("BN_API_SECRET")
    mac := hmac.New(sha256.New, []byte(api_secret))
    mac.Write([]byte(totalParams))
    sign := hex.EncodeToString(mac.Sum(nil))
    values.Add("signature", sign)
  }

  // prepare request
  req, _ := http.NewRequest(method, base_url + path + "?" + values.Encode(), nil)

  //set header
  api_key := os.Getenv("BN_API_KEY")
  req.Header.Add("X-MBX-APIKEY", api_key)
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  // do request
  client := new(http.Client)
  resp, err = client.Do(req)
  return
}
