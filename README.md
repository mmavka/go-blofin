# go-blofin

Библиотека для работы с криптобиржей BloFin на Go

## Описание

`go-blofin` — это модульная библиотека для работы с REST и WebSocket API BloFin, вдохновлённая архитектурой go-binance. Поддерживает публичные и приватные методы, подписку на события, генерацию подписей, обработку ошибок и покрыта unit-тестами.

## Возможности
- Публичные REST методы: инструменты, тикеры, стакан, сделки, свечи, mark price, funding rate
- Приватные REST методы: баланс, позиции, ордера, история переводов и выводов
- WebSocket: подписка на trades, candles, tickers, order book, funding-rate
- Генерация подписей для приватных запросов
- Полное покрытие unit-тестами
- Современная архитектура, расширяемость, чистый код

## Структура проекта
- `rest/` — REST-клиент, сервисы, модели
- `ws/` — WebSocket-клиент, модели, подписи
- `auth/` — генерация подписей
- `utils/` — базовые ошибки
- `docs/` — документация (архитектура, changelog, задачи, Q&A)

## Быстрый старт
```go
import "github.com/mmavka/go-blofin/rest"

client := rest.NewDefaultRestClient()
resp, err := client.NewGetInstrumentsService().Do(context.Background())
if err != nil {
    // обработка ошибки
}
for _, inst := range resp.Data {
    fmt.Println(inst.InstID)
}
```

## Пример WebSocket
```go
import "github.com/mmavka/go-blofin/ws"

wsClient := ws.NewDefaultClient()
wsClient.SetErrorHandler(func(err error) {
    log.Printf("WebSocket error: %v", err)
})
if err := wsClient.Connect(); err != nil {
    panic(err)
}
_ = wsClient.Subscribe([]ws.ChannelArgs{{Channel: "trades", InstId: "BTC-USDT"}})
for trade := range wsClient.Trades() {
    fmt.Println(trade)
}
```

## Обработка ошибок WebSocket
```go
wsClient := ws.NewDefaultClient()
wsClient.SetErrorHandler(func(err error) {
    log.Printf("WebSocket error: %v", err)
})
```

## Кастомный endpoint
```go
import "github.com/mmavka/go-blofin/rest"
import "github.com/mmavka/go-blofin/ws"

client := rest.NewDefaultRestClient()
client.SetBaseURL("https://sandbox.blofin.com")

wsClient := ws.NewDefaultClient()
wsClient.SetURL("wss://sandbox-ws.blofin.com/ws")
```

## Особенности подписки на WebSocket-каналы
- Публичные каналы (trades, candles, tickers, order book и др.) не требуют аутентификации.
- Приватные каналы (orders, positions, orders-algo и др.) требуют предварительного вызова Login.
- Общая длина запроса на подписку не должна превышать 4096 байт.
- При нарушении этих условий будет возвращена ошибка.

## Требования
- Go 1.24+
- [resty](https://github.com/go-resty/resty)
- [goccy/go-json](https://github.com/goccy/go-json)
- [gorilla/websocket](https://github.com/gorilla/websocket)

## Лицензия
MIT 