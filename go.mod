module github.com/egelis/binance

go 1.17

require (
	github.com/adshao/go-binance/v2 v2.3.1
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/shopspring/decimal v1.2.0
)

require (
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/klauspost/compress v1.13.1 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

replace github.com/adshao/go-binance/v2 => github.com/egelis/go-binance/v2 v2.3.2
