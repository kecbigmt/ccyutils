package bitbank

import (
  "fmt"
  //"testing"
)

func ExampleMarkets() {
  pairs, _ := AvailablePair()
  b := pairs.Have("ETH_BTC", "bitbank")
  fmt.Println(b)
  // Output: true
}
