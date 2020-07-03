package bitflyer

import (
  "fmt"
  //"testing"
)

func ExampleMarkets() {
  pairs, _ := AvailablePair()
  b := pairs.Have("ETH_BTC", "bitflyer")
  fmt.Println(b)
  // Output: true
}
