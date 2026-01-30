<div align="center">

# SIEM-Agent

**Security Event Collector**

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)

<p align="center">
Агент для сбора, нормализации и отправки событий безопасности в NoSQLdb
</p>

---

</div>

## Особенности

- **Мониторинг логов** — отслеживание изменений в реальном времени через fsnotify
- **Парсеры** — поддержка syslog, auth.log, bash history
- **Фильтрация** — гибкая настройка include/exclude паттернов
- **Буферизация** — сохранение событий на диск при недоступности БД
- **Батчинг** — отправка событий пачками для оптимизации

---

## Поддерживаемые источники

| Источник | Описание |
|----------|----------|
| `/var/log/syslog` | Системные события |
| `/var/log/auth.log` | События аутентификации |
| `/var/log/messages` | Общие системные сообщения |
| `~/.bash_history` | История команд пользователя |

---

## Быстрый старт

### Docker (рекомендуется)

Агент запускается автоматически через docker-compose вместе со всем стеком:

```bash
docker compose up --build
```

### Локальный запуск

```bash
# Сборка
go build -o siem-agent ./cmd/agent/main.go

# Запуск
CONFIG_PATH=./configs/config.yaml ./siem-agent
```

---

## Конфигурация

### Основные параметры

```yaml
agent_id: "agent-001"
db_address: "localhost:5140"
db_name: "siem_events"

buffer:
  path: "/tmp/siem-agent-buffer"
  max_size: 10485760  # 10 MB

batch:
  size: 100
  interval: 5s

log_sources:
  - path: "/var/log/syslog"
    type: "syslog"
  - path: "/var/log/auth.log"
    type: "auth"

filter:
  severity_threshold: "low"
  exclude_patterns:
    - "CRON.*session"
  include_patterns:
    - "ssh|sudo|auth"
```

### Файлы конфигурации

| Файл | Описание |
|------|----------|
| [configs/config.yaml](configs/config.yaml) | Локальная конфигурация |
| [configs/config.docker.yaml](configs/config.docker.yaml) | Конфигурация для Docker |

---

## Архитектура

```
SIEM-Agent/
├── cmd/agent/              # Точка входа
├── configs/                # Конфигурационные файлы
├── deployments/systemd/    # Systemd unit файл
├── internal/
│   ├── config/             # Загрузка конфигурации
│   ├── domain/             # Модели данных (Event, Batch)
│   ├── filter/             # Фильтрация событий
│   ├── logger/             # Логирование
│   ├── parser/             # Парсеры логов
│   ├── reader/             # Чтение файлов (fsnotify)
│   ├── sender/             # Отправка в БД
│   └── storage/            # Дисковый буфер
└── scripts/                # Установка и удаление
```

---

## Типы событий

| Тип | Severity | Описание |
|-----|----------|----------|
| `user_login` | medium | Успешный вход пользователя |
| `user_logout` | low | Выход пользователя |
| `auth_failure` | high | Неудачная попытка входа |
| `sudo_command` | medium | Выполнение команды через sudo |
| `system_event` | info | Общее системное событие |

---

## Тестирование

```bash
# Unit-тесты
go test ./...

# Тесты с детектором гонок
go test -race ./...
```
