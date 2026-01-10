DB_DIR=NoSQLdb
AGENT_DIR=SIEM-Agent
WEB_DIR=Web/backend
BIN_DIR=bin

BLUE=\033[0;34m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m

.PHONY: all build run-db run-agent run-web run-all stop clean help

help:
	@echo -e "$(BLUE)Доступные команды:$(NC)"
	@echo "  make run-db     - Запустить NoSQLdb"
	@echo "  make run-web    - Запустить Web-Backend"
	@echo "  make run-agent  - Запустить SIEM-Agent (требуется sudo)"
	@echo "  make stop       - Остановить все запущенные компоненты"

run-db:
	@echo -e "$(GREEN)Запуск NoSQLdb на порту 5140...$(NC)"
	cd $(DB_DIR) && go run ./cmd/server/main.go

run-web:
	@echo -e "$(GREEN)Запуск Web-Backend на порту 8080...$(NC)"
	cd $(WEB_DIR) && go run ./cmd/server/main.go

run-agent:
	@echo -e "$(GREEN)Запуск SIEM-Agent...$(NC)"
	cd $(AGENT_DIR) && sudo go run ./cmd/agent/main.go


stop:
	@echo -e "$(YELLOW)Остановка всех компонентов...$(NC)"
	@if [ -f /tmp/siem_agent.pid ]; then sudo kill $$(cat /tmp/siem_agent.pid) 2>/dev/null || true; rm -f /tmp/siem_agent.pid; echo "SIEM-Agent остановлен"; fi
	@if [ -f /tmp/web_backend.pid ]; then kill $$(cat /tmp/web_backend.pid) 2>/dev/null || true; rm -f /tmp/web_backend.pid; echo "Web-Backend остановлен"; fi
	@if [ -f /tmp/nosql_db.pid ]; then kill $$(cat /tmp/nosql_db.pid) 2>/dev/null || true; rm -f /tmp/nosql_db.pid; echo "NoSQLdb остановлен"; fi
	@pkill -f "go run.*cmd/server/main.go" 2>/dev/null || true
	@pkill -f "go run.*cmd/agent/main.go" 2>/dev/null || true
	@echo -e "$(BLUE)Все компоненты остановлены$(NC)"
