package repository

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Narotan/Web-SIEM/Web/backend/internal/repository/model"
)

type Repository interface {
	FindAll(database string, query map[string]any) ([]map[string]any, error)
}

type nosqlRepository struct {
	addr string
}

func NewNosqlRepository(addr string) Repository {
	return &nosqlRepository{
		addr: addr,
	}
}

func (r *nosqlRepository) FindAll(database string, query map[string]any) ([]map[string]any, error) {
	req := model.DBRequest{
		Database: database,
		Command:  "find",
		Query:    query,
	}

	conn, err := net.DialTimeout("tcp", r.addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к СУБД по адресу %s: %w", r.addr, err)
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(req); err != nil {
		return nil, fmt.Errorf("ошибка кодирования запроса: %w", err)
	}

	var resp model.DBResponse
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа от СУБД: %w", err)
	}

	if resp.Status == "error" {
		return nil, fmt.Errorf("%s", resp.Message)
	}

	return resp.Data, nil
}
