<div align="center">

# Web-SIEM

**Security Information and Event Management**

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

<p align="center">
  <b>Учебный SIEM-стенд с полным циклом: сбор логов, хранение, анализ и визуализация</b>
</p>

---

</div>

## Обзор

Web-SIEM — это учебный проект, демонстрирующий архитектуру системы управления событиями безопасности. Агент собирает логи с хостов, нормализует их и отправляет в документную NoSQL-базу данных. Веб-панель отображает статистику, позволяет фильтровать события и экспортировать данные.

### Ключевые особенности

| Компонент | Описание |
|-----------|----------|
| **NoSQLdb** | Собственная документная СУБД с B+Tree индексами и TCP-интерфейсом |
| **SIEM-Agent** | Сборщик логов с поддержкой syslog, auth.log, bash history |
| **Web API** | REST API на Gin с кэшированием статистики и пагинацией |
| **Dashboard** | Интерактивная панель с графиками и экспортом в JSON/CSV |

---

## Архитектура

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  SIEM-Agent │───▶│   NoSQLdb   |◀── │   Web API   │
│  (collector)│ TCP │  (storage)  │ TCP │   (REST)    │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                    ┌─────────────┐            │
                    │    nginx    │◀──────────┘
                    │   (proxy)   │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  Dashboard  │
                    │    (UI)     │
                    └─────────────┘
```

---

## Быстрый старт

### Требования

- Docker и Docker Compose
- Git

### Запуск

```bash
# Клонирование репозитория
git clone https://github.com/shushuh3/Web-SIEM.git
cd Web-SIEM

# Создание файла окружения (опционально)
cp .env.example .env

# Запуск всех сервисов
docker compose up --build
```

### Доступ

| Сервис | URL | Порт |
|--------|-----|------|
| Dashboard | http://localhost | 80 |
| Web API | http://localhost:8080/api | 8080 |
| NoSQLdb | TCP | 5140 |

> **Авторизация:** По умолчанию используется BasicAuth с логином `admin` и паролем `admin`

---

## Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_SOCKET` | Адрес NoSQLdb | `nosql-db:5140` |
| `DB_NAME` | Имя базы данных | `siem_events` |
| `SERVER_PORT` | Порт Web API | `8080` |
| `WEB_USER` | Логин для авторизации | `admin` |
| `WEB_PASSWORD` | Пароль для авторизации | `admin` |

---

## API

### Эндпойнты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/health` | Проверка состояния сервиса |
| GET | `/api/events` | Получение событий с пагинацией |
| GET | `/api/events/export` | Экспорт событий в JSON или CSV |
| GET | `/api/stats` | Статистика для дашборда |

### Примеры запросов

```bash
# Проверка здоровья
curl -u admin:admin http://localhost:8080/api/health

# Получение событий
curl -u admin:admin "http://localhost:8080/api/events?page=1&limit=50"

# Экспорт в CSV
curl -u admin:admin "http://localhost:8080/api/events/export?format=csv" -o events.csv
```

---

## Структура проекта

```
Web-SIEM/
├── NoSQLdb/                # Документная база данных
│   ├── cmd/                # Точки входа (сервер, клиент)
│   ├── internal/           # Внутренняя логика
│   │   ├── handlers/       # Обработчики команд
│   │   ├── index/          # B+Tree индексы
│   │   ├── operators/      # Операторы сравнения
│   │   ├── query/          # Парсер запросов
│   │   └── storage/        # Коллекции и персистентность
│   └── tests/              # Интеграционные тесты
│
├── SIEM-Agent/             # Агент сбора логов
│   ├── cmd/agent/          # Точка входа
│   ├── configs/            # Конфигурационные файлы
│   └── internal/           # Внутренняя логика
│       ├── filter/         # Фильтрация событий
│       ├── parser/         # Парсеры логов
│       ├── reader/         # Чтение файлов
│       └── sender/         # Отправка в БД
│
├── Web/
│   ├── backend/            # REST API
│   │   ├── cmd/server/     # Точка входа
│   │   └── internal/       # Сервисы, репозитории, транспорт
│   └── frontend/           # Статический UI
│       ├── css/            # Стили
│       └── js/             # JavaScript модули
│
├── nginx/                  # Конфигурация прокси
├── docker-compose.yml      # Оркестрация контейнеров
└── .env.example            # Пример переменных окружения
```

---

## Тестирование

```bash
# Тесты NoSQLdb
cd NoSQLdb && go test ./...

# Тесты Web API
cd Web/backend && go test ./...

# Тесты SIEM-Agent
cd SIEM-Agent && go test ./...

# Запуск с детектором гонок
go test -race ./...
```

---

## Технологии

- **Backend:** Go 1.23+, Gin Framework
- **Frontend:** Vanilla JS, Chart.js
- **Database:** Собственная NoSQL с B+Tree индексами
- **Infrastructure:** Docker, Docker Compose, nginx

---
