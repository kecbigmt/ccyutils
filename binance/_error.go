import binance


type Error struct{
  ErrorCode string `json:"error_code"`
  Component string `json:"component"`
  Service string `json:"service"`
  Message string `json:""`
  "error_code":"500",
  "component":"order",
  "service":"binance",
  "message":"[Error]Failed to request",
  "detail":%v
}
