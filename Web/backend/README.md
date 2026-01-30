<div align="center">

# Web Backend

**SIEM REST API**

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge)](https://gin-gonic.com/)

<p align="center">
REST API для получения событий безопасности, статистики и экспорта данных
</p>

---

</div>

## Особенности

- **REST API** — стандартные HTTP эндпойнты
- **BasicAuth** — защита всех эндпойнтов
- **Кэширование** — TTL-кэш статистики для снижения нагрузки
- **Пагинация** — постраничный вывод событий
- **Экспорт** — выгрузка в JSON и CSV форматы

---

## API Reference

### Эндпойнты

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/health` | Проверка состояния сервиса |
| `GET` | `/api/events` | Получение событий с пагинацией |
| `GET` | `/api/events/export` | Экспорт всех событий |
| `GET` | `/api/stats` | Статистика для дашборда |

### Параметры запросов

#### GET /api/events

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|--------------|----------|
| `page` | int | 1 | Номер страницы |
| `limit` | int | 50 | Количество записей (макс. 200) |

#### GET /api/events/export

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|--------------|----------|
| `format` | string | json | Формат экспорта: `json` или `csv` |

---

## Примеры

```bash
# Health check
curl -u admin:admin http://localhost:8080/api/health

# Получение событий
curl -u admin:admin "http://localhost:8080/api/events?page=1&limit=50"

# Статистика
curl -u admin:admin http://localhost:8080/api/stats

# Экспорт в CSV
curl -u admin:admin "http://localhost:8080/api/events/export?format=csv" -o events.csv

# Экспорт в JSON
curl -u admin:admin "http://localhost:8080/api/events/export?format=json" -o events.json
```

---

## Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_SOCKET` | Адрес NoSQLdb | `localhost:5140` |
| `DB_NAME` | Имя базы данных | `siem_events` |
| `SERVER_PORT` | Порт сервера | `8080` |
| `WEB_USER` | Логин BasicAuth | `admin` |
| `WEB_PASSWORD` | Пароль BasicAuth | `admin` |

---

## Архитектура

```
backend/
├── cmd/server/             # Точка входа
└── internal/
    ├── config/             # Загрузка конфигурации
    ├── repository/         # Слой доступа к данным
    │   └── model/          # Модели запросов к БД
    ├── service/            # Бизнес-логика
    │   └── domain/         # Доменные модели
    └── transport/
        ├── dto/            # Data Transfer Objects
        └── http/           # HTTP handlers и роутинг
            └── middleware/ # CORS, BasicAuth
```

---

## Запуск

### Docker (рекомендуется)

```bash
docker compose up --build web-backend
```

### Локально

```bash
# Установка зависимостей
go mod download

# Запуск
go run ./cmd/server/main.go
```

---

## Тестирование

```bash
# Unit-тесты
go test ./...

# С подробным выводом
go test -v ./...

# С детектором гонок
go test -race ./...
```
