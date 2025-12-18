# Lab3 - Nobel Prize REST API

REST API для работы с данными о Нобелевских лауреатах и премиях.

## Требования

- Go 1.22+ (проект использует go 1.25.1)
- PostgreSQL 14+

## Конфигурация

Переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | URL подключения к PostgreSQL | `postgres://postgres:postgres@localhost:5432/ris` |
| `PORT` | Порт для HTTP сервера | `8080` |
| `API_TOKEN` | Токен для авторизации | `secret-api-token` |

## Запуск

```bash
# Установка зависимостей
go mod download

# Запуск сервера
go run ./cmd/lab3

# Или сборка и запуск
go build -o lab3 ./cmd/lab3
./lab3
```

## API Endpoints

### Аутентификация

API требует авторизации одним из способов:

1. **Bearer Token** в заголовке:
   ```
   Authorization: Bearer secret-api-token
   ```

2. **API Key** в параметре запроса:
   ```
   ?api_key=secret-api-token
   ```

### Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/health` | Проверка здоровья сервиса (без авторизации) |
| GET | `/swagger/` | Swagger UI документация |
| GET | `/api/v1/stats` | Статистика набора данных |
| GET | `/api/v1/stats/last-update` | Дата последнего обновления |
| GET | `/api/v1/categories` | Список категорий премий |
| GET | `/api/v1/laureates` | Список лауреатов (с пагинацией) |
| GET | `/api/v1/laureates/:id` | Получить лауреата по ID |
| POST | `/api/v1/laureates` | Создать лауреата |
| PUT | `/api/v1/laureates/:id` | Обновить лауреата |
| DELETE | `/api/v1/laureates/:id` | Удалить лауреата |
| GET | `/api/v1/prizes` | Список премий (с пагинацией) |
| GET | `/api/v1/prizes/:id` | Получить премию по ID |
| GET | `/api/v1/prizes/category/:category` | Премии по категории |
| GET | `/api/v1/prizes/year/:year` | Премии по году |
| POST | `/api/v1/prizes` | Создать премию |
| PUT | `/api/v1/prizes/:id` | Обновить премию |
| DELETE | `/api/v1/prizes/:id` | Удалить премию |

## Примеры запросов

### Получить статистику
```bash
curl -H "Authorization: Bearer secret-api-token" http://localhost:8080/api/v1/stats
```

### Получить список лауреатов
```bash
curl -H "Authorization: Bearer secret-api-token" "http://localhost:8080/api/v1/laureates?page=1&per_page=10"
```

### Создать лауреата
```bash
curl -X POST -H "Authorization: Bearer secret-api-token" \
     -H "Content-Type: application/json" \
     -d '{"id": 999, "firstname": "Test", "surname": "User", "motivation": "For testing", "share": 1}' \
     http://localhost:8080/api/v1/laureates
```

### Получить премии по категории
```bash
curl -H "Authorization: Bearer secret-api-token" http://localhost:8080/api/v1/prizes/category/physics
```

### С использованием API Key
```bash
curl "http://localhost:8080/api/v1/stats?api_key=secret-api-token"
```

## Swagger UI

Документация API доступна по адресу: `http://localhost:8080/swagger/`

## Структура проекта

```
cmd/lab3/
├── main.go          # Точка входа приложения
└── swagger.go       # Swagger спецификация

internal/app/api/
├── middleware/
│   └── auth.go      # Middleware авторизации
└── v1/
    ├── dto.go       # Data Transfer Objects
    ├── handlers.go  # HTTP handlers
    ├── routes.go    # Регистрация маршрутов
    └── service.go   # Бизнес-логика

pkg/postgres/queries/
├── laureates.sql    # SQL запросы для лауреатов
├── prizes.sql       # SQL запросы для премий
└── stats.sql        # SQL запросы для статистики
```

## База данных

### Схема

```sql
CREATE TABLE laureates (
    id INT PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    surname VARCHAR(100),
    motivation TEXT NOT NULL,
    share INT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE prizes (
    id SERIAL PRIMARY KEY,
    year INT NOT NULL,
    category VARCHAR(100) NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE prizes_to_laureates (
    prize_id INT REFERENCES prizes(id) ON DELETE CASCADE,
    laureate_id INT REFERENCES laureates(id) ON DELETE CASCADE,
    PRIMARY KEY (prize_id, laureate_id)
);
```

### Применение миграций

```sql
-- Добавление колонки updated_at если её нет
ALTER TABLE laureates ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();
ALTER TABLE prizes ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();
```
