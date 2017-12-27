package bitbank

import (
  "io"
  "os"
  "log"
  "fmt"
  "time"
  "strconv"
  "strings"
  "net/http"
  "net/url"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"
)

func AuthorizedRequest(method string, path string, params map[string]string) (resp *http.Response, err error){
  // init
  base_url := "https://api.bitbank.cc/v1"
  timestamp := strconv.FormatInt(time.Now().Unix(), 10)

  // prepare request
  var body io.Reader
  var content string
  switch method {
  case "GET":
    if len(params)>0{
      values := url.Values{}
      for key, value := range params{
        values.Add(key, value)
      }
      path = path + "?" + values.Encode()
    }
    body  = nil
    content = timestamp + "/v1" +path
  case "POST":
    jsonBytes, err := json.Marshal(params)
    if err != nil {
        log.Fatal(err)
    }
    body = strings.NewReader(string(jsonBytes))
    content = timestamp + string(jsonBytes)
  default:
    log.Fatal(fmt.Sprintf("Invalid HTTP method: %v", method))
  }
  req, _ := http.NewRequest(method, base_url + path, body)

  //add authorization info
  api_secret := os.Getenv("BB_API_SECRET")
  api_key := os.Getenv("BB_API_KEY")
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("ACCESS-KEY", api_key)
  req.Header.Add("ACCESS-NONCE", timestamp)
  // HMAC-SHA256 sign
  mac := hmac.New(sha256.New, []byte(api_secret))
  mac.Write([]byte(content))
  sign := hex.EncodeToString(mac.Sum(nil))
  req.Header.Add("ACCESS-SIGNATURE", sign)

  // do request
  client := new(http.Client)
  resp, err = client.Do(req)
  return
}
