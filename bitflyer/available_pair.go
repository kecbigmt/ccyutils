package bitflyer

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/kecbigmt/ccyutils"
)

// struct to receive bitFlyer format JSON
type BitFlyerAvailablePairRaw struct {
  ProductCode string `json:"product_code"`
}

// get available currency pairs
func AvailablePair() (output ccyutils.AvailablePairs){
    url := "https://api.bitflyer.jp/v1/markets"
    req, _ := http.NewRequest("GET", url, nil)
    client := new(http.Client)
    resp, err := client.Do(req)
    if err != nil{
      err = fmt.Errorf(`{
        "error_code":"800",
        "component":"markets",
        "service":"bitflyer",
        "message":"[Error]Failed to request",
        "detail":%v
      }`, err)
      log.Fatal(err)
    }
    defer resp.Body.Close()
    bytes, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200{
      err = fmt.Errorf(`{
        "error_code":"801",
        "component":"markets",
        "service":"bitflyer",
        "message":"[Error]HTTP Error(%v)",
        "detail":%v
      }`, resp.StatusCode, string(bytes))
      log.Fatal(err)
    }
    var json_receiver []BitFlyerAvailablePairRaw
    if err = json.Unmarshal(bytes, &json_receiver); err != nil {
      err = fmt.Errorf(`{
        "error_code":"802",
        "component":"markets",
        "service":"bitflyer",
        "message":"[Error]Failed to decode JSON",
        "detail":%v
      }`, err)
      log.Fatal(err)
    }
    var sl []string
    for _, v := range json_receiver{
      sl = append(sl, v.ProductCode)
    }
    output = ccyutils.AvailablePairs(sl)
    return
}
