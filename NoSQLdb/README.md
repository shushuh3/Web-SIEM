<div align="center">

# NoSQLdb

**Lightweight Document-Oriented Database**

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)

<p align="center">
Учебная документно-ориентированная NoSQL СУБД с сетевым интерфейсом и поддержкой конкурентных write-операций
</p>

---

</div>

## Особенности

- **Документная модель** — хранение коллекций JSON-документов
- **TCP-сервер** — клиент-серверная архитектура, работа по сети
- **REPL-клиент** — интерактивный режим командной строки
- **B+Tree индексы** — быстрый поиск по индексированным полям
- **Гибкие запросы** — операторы `$eq`, `$gt`, `$lt`, `$in`, `$like`, `$or`, `$and`
- **Очередь write-операций** — гарантированная последовательность изменений
- **Потокобезопасность** — конкурентный доступ к коллекциям
- **Персистентность** — хранение данных и индексов на диске

> В проекте используются собственные реализации B+Tree и HashMap

---

## Быстрый старт

### Запуск сервера

```bash
go run ./cmd/server/main.go
```

### Запуск клиента

```bash
go run ./cmd/client/main.go --host localhost --port 5140 --database my_database
```

---

## Примеры команд

```sql
-- Вставка документа
INSERT users {"name": "Alice", "age": 25}

-- Поиск с условием
FIND users {"age": {"$gt": 20}}

-- Удаление документа
DELETE users {"name": "Alice"}

-- Создание индекса
CREATE_INDEX users age
```

---

## Архитектура

```
NoSQLdb/
├── cmd/
│   ├── server/         # TCP-сервер
│   └── client/         # REPL-клиент
├── internal/
│   ├── handlers/       # Обработчики команд (INSERT, FIND, DELETE)
│   ├── index/          # B+Tree индексы
│   ├── operators/      # Операторы сравнения ($eq, $gt, $lt, $like)
│   ├── query/          # Парсер JSON-запросов
│   ├── server/         # TCP-сервер и роутинг
│   └── storage/        # Коллекции, HashMap, менеджер, персистентность
└── tests/              # Интеграционные тесты конкурентности
```

---

## Очередь задач

Все операции изменения (insert, delete, create_index) ставятся в очередь. Воркер по одной обрабатывает задачи, гарантируя целостность данных. Результат возвращается через канал обратно вызывающему хендлеру.

Подробнее: [internal/storage/manager.go](internal/storage/manager.go)

---

## Тестирование

```bash
# Запуск сервера
go run ./cmd/server/main.go

# В другом терминале — тесты
go test ./... -v

# Тесты конкурентности (требуют запущенный сервер)
go test -v ./tests/
```
