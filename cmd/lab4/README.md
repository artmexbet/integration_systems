# Lab 4 - NATS Event Streaming

Этот модуль содержит два приложения для работы с событиями через NATS JetStream:

## 1. Event Listener (`main.go`)

Основной сервис, который слушает события из стрима и выводит их в консоль.

### Использование:
```bash
go run ./cmd/lab4
```

### Функционал:
- Подключается к NATS серверу (по умолчанию `nats://localhost:4222`)
- Создаёт поток событий `EVENTS` с двумя субъектами:
  - `prize.created` - события о созданных премиях
  - `laureate.created` - события о созданных лауреатах
- Слушает события и выводит их в формате JSON

### Пример вывода:
```
Listening for events from stream...
Received prize created event prize={
  "year": "2023",
  "category": "Physics",
  ...
}
```

## 2. Get Last Message (`cmd/lab4/get-last-msg/main.go`)

Утилита для получения последнего сообщения из стрима.

### Использование:
```bash
# Получить последнее сообщение о премии
go run ./cmd/lab4/get-last-msg -type=prize

# Получить последнее сообщение о лауреате
go run ./cmd/lab4/get-last-msg -type=laureate
```

### Флаги:
- `-type` (string, default="prize") - тип сообщения: `prize` или `laureate`

### Пример вывода:
```
Last Prize Message:
Year: 2023
Category: Physics
Overall Motivation: For groundbreaking contributions...
Number of Laureates: 3
  Laureate 1: John Smith
  Laureate 2: Jane Doe
  Laureate 3: Bob Johnson
```

## API

### Subscriber Interface

Пакет `internal/subscriber` предоставляет следующие методы:

```go
// Подписаться на события о новых премиях
SubscribePrizeCreated(handler func(prize domain.Prize) error) error

// Подписаться на события о новых лауреатах
SubscribeLaureateCreated(handler func(laureate domain.Laureate) error) error

// Получить последнее сообщение о премии из стрима
GetLastPrizeMessage() (*domain.Prize, error)

// Получить последнее сообщение о лауреате из стрима
GetLastLaureateMessage() (*domain.Laureate, error)

// Закрыть соединение и отписаться от всех событий
Close() error
```

## Структура потока

Поток NATS JetStream `EVENTS` настроен с следующими параметрами:
- **Срок хранения**: 7 дней
- **Политика подтверждения**: Явное подтверждение (AckExplicit)
- **Политика доставки**: Доставить все накопленные сообщения (DeliverAll)

