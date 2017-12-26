# About
いろんな仮想通貨取引所のPublic APIやPrivate APIを叩きます。
現在対応しているのは、
* bitFlyer
   * Ticker
   * Balance（アカウント残高取得）
* bitbank.cc
  * Ticker
* Binance
  * Ticker

※随時対応APIを増やしたい
# Getting Started
```
go get github.com/kecbigmt/ccyutils
```
# Example
## Tickerを取得する
1. `ccyutils/{取引所の名前}`でパッケージ読み込み
2. `Ticker({通貨ペア})`関数で取得
```
package main

import (
  "fmt"
  bf "github.com/kecbigmt/ccyutils/bitflyer"
  )

func main(){
  ticker, _ := bf.Ticker("BTC_JPY")
  fmt.Println(ticker)
}
```
