# Lab 3: SOAP Web Service Implementation

Реализация SOAP веб-сервиса для управления файлами с поддержкой WSDL.

## Описание

Этот проект реализует SOAP веб-сервис для загрузки и управления файлами с полной поддержкой WSDL, аутентификацией и асинхронными уведомлениями.

## Структура проекта

```
cmd/lab3/
├── README.md           # Задание (исходное)
├── IMPLEMENTATION.md   # Этот файл с инструкциями
├── main.go            # Точка входа с описанием
├── server/
│   └── main.go        # SOAP сервер
├── client/
│   └── main.go        # Интерактивный клиент
└── test_server.sh     # Скрипт автоматического тестирования
```

## Возможности сервера

### 1. WSDL описание сервиса
- Автоматическая генерация WSDL
- Доступ: `http://localhost:8080/soap?wsdl`

### 2. Аутентификация
- Базовая аутентификация через SOAP Header
- Учетные данные:
  - `user1:pass1`
  - `user2:pass2`
  - `admin:admin`

### 3. Методы сервиса

#### UploadFile
Загрузка файла на сервер с валидацией:
- Максимальный размер: 3 МБ
- Кодирование: Base64 (MTOM)
- Асинхронное уведомление клиента о результате
- Валидация:
  - Проверка на пустой файл
  - Проверка размера (до 3 МБ)
  - Запрет файлов с буквой 'Ж' в имени
  - Запрет файлов, содержащих только JSON
  - Проверка доступного места для хранения

#### GetLastFileInfo
Получение информации о последнем загруженном файле текущего пользователя:
- Имя файла
- Размер
- Время загрузки

#### GetFileListCSV
Получение списка всех файлов на сервере в формате CSV:
- Данные всех пользователей
- Формат: Username, FileName, FileSize, UploadTime
- Кодирование: Base64

#### GetUptime
Получение времени работы сервера.

### 4. Асинхронные уведомления
- Отправка результата загрузки файла на указанный callback URL
- Timeout: 10 секунд

## Запуск сервера

```bash
# Из корня проекта
go run cmd/lab3/server/main.go

# Сервер запустится на порту 8080
# WSDL: http://localhost:8080/soap?wsdl
```

## Запуск клиента

```bash
# Из корня проекта
go run cmd/lab3/client/main.go
```

### Возможности клиента

1. **Проверка доступности сервера** - автоматическая при запуске
2. **Аутентификация** - запрос учетных данных
3. **Загрузка файла** - выбор файла и отправка с ожиданием уведомления
4. **Проверка последнего файла** - информация о последнем загруженном файле
5. **Получение списка файлов** - CSV список всех файлов на сервере
6. **Получение uptime сервера**
7. **Webhook сервер** - автоматический запуск на порту 9090 для получения уведомлений
8. **Проактивная проверка** - если уведомление не получено в течение 15 секунд

## Автоматическое тестирование

```bash
# Запустите сервер в одном терминале
go run cmd/lab3/server/main.go

# В другом терминале запустите тесты
bash cmd/lab3/test_server.sh
```

Скрипт проверит:
- ✓ Доступность сервера и WSDL
- ✓ Метод GetUptime
- ✓ Загрузку файла
- ✓ Получение информации о последнем файле
- ✓ Получение списка файлов (CSV)
- ✓ Валидацию: запрет буквы 'Ж' в имени
- ✓ Валидацию: запрет JSON файлов
- ✓ Валидацию: запрет пустых файлов
- ✓ Проверку аутентификации

## Примеры SOAP запросов

### GetUptime

```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetUptime xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
```

### UploadFile

```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>test.txt</fileName>
      <fileData>VGVzdCBkYXRh</fileData>
      <callbackURL>http://localhost:9090/webhook</callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
```

### GetLastFileInfo

```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetLastFileInfo xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
```

### GetFileListCSV

```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetFileListCSV xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
```

## Отправка запросов через curl

```bash
# Получить WSDL
curl http://localhost:8080/soap?wsdl

# Отправить SOAP запрос
curl -X POST http://localhost:8080/soap \
  -H "Content-Type: text/xml" \
  -d @request.xml
```

## Технические детали

### Использованные технологии
- Go 1.25.1
- Стандартная библиотека Go (encoding/xml, net/http)
- SOAP 1.1 / WSDL 1.1

### Хранилище файлов
- Директория: `./uploaded_files/`
- Максимальный размер хранилища: 100 МБ
- Имя файла: `{username}_{filename}`

### Безопасность
- Базовая аутентификация через SOAP Header
- Валидация всех входных данных
- Ограничение размера файлов
- Защита от загрузки потенциально опасных файлов

## Соответствие требованиям задания

### Требования к серверу
- [x] Предоставление WSDL описания
- [x] Аутентификация
- [x] Загрузка файлов до 3 МБ с MTOM
- [x] Асинхронное уведомление клиента
- [x] Метод получения информации о последнем файле
- [x] Метод получения списка файлов (CSV)
- [x] Метод получения uptime
- [x] Валидация файлов (все 5 правил)

### Требования к клиенту
- [x] Проверка доступности сервера
- [x] Аутентификация
- [x] Выбор файла для отправки
- [x] Webhook сервер для уведомлений
- [x] Проактивная проверка (fallback)
- [x] Отображение списка файлов

## Расширения и улучшения

Возможные улучшения (не входят в базовое задание):
- Персистентное хранилище (база данных)
- TLS/SSL для безопасного соединения
- Более сложная схема аутентификации (JWT, OAuth)
- Ограничение скорости (rate limiting)
- Логирование в файл
- Метрики и мониторинг
- Docker контейнеризация
- Модульные тесты

## Troubleshooting

### Сервер не запускается
- Проверьте, не занят ли порт 8080: `lsof -i :8080`
- Убедитесь, что Go установлен: `go version`

### Клиент не может подключиться
- Убедитесь, что сервер запущен
- Проверьте доступность: `curl http://localhost:8080/soap?wsdl`

### Уведомления не приходят
- Это нормально, если клиент не запущен или webhook сервер недоступен
- Клиент автоматически использует проактивную проверку через 15 секунд
