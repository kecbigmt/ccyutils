package bitbank

import (
  "github.com/kecbigmt/ccyutils"
)

// get available currency pairs *bitbank has no API to get available pairs.
func AvailablePair() (output ccyutils.AvailablePairs){
  sl := []string{
    "btc_jpy",
    "xrp_jpy",
    "ltc_btc",
    "eth_btc",
    "mona_jpy",
    "mona_btc",
    "bcc_jpy",
    "bcc_btc",
  }
  output = ccyutils.AvailablePairs(sl)
  return
}
