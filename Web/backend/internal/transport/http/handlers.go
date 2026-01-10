package http

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Narotan/Web-SIEM/Web/backend/internal/service"
	"github.com/Narotan/Web-SIEM/Web/backend/internal/transport/dto"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:  "ok",
		Message: "SIEM Web API работает",
	})
}

func (h *Handler) GetEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	result, err := h.service.GetEvents(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status: "error",
			Error:  "Ошибка связи с базой данных: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.EventsResponse{
		Status:     "success",
		Count:      result.Count,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
		Data:       result.Data,
	})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.StatsResponse{
		ActiveAgents:  stats.ActiveAgents,
		EventsByType:  stats.EventsByType,
		SeverityDist:  stats.SeverityDist,
		TopUsers:      stats.TopUsers,
		TopProcesses:  stats.TopProcesses,
		EventsPerHour: stats.EventsPerHour,
		LastLogins:    stats.LastLogins,
	})
}

func (h *Handler) ExportEvents(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	events, err := h.service.ExportAllEvents()
	if err != nil {
		log.Printf("Export error: %v", err) // Log error to console
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status: "error",
			Error:  "Ошибка экспорта данных: " + err.Error(),
		})
		return
	}

	timestamp := time.Now().Format("20060102_150405")

	switch format {
	case "csv":
		h.exportCSV(c, events, timestamp)
	case "json":
		h.exportJSON(c, events, timestamp)
	default:
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status: "error",
			Error:  "Неподдерживаемый формат. Используйте 'json' или 'csv'",
		})
	}
}

func (h *Handler) exportJSON(c *gin.Context, events []map[string]any, timestamp string) {
	filename := fmt.Sprintf("events_export_%s.json", timestamp)
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status: "error",
			Error:  "Ошибка формирования JSON: " + err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func (h *Handler) exportCSV(c *gin.Context, events []map[string]any, timestamp string) {
	filename := fmt.Sprintf("events_export_%s.csv", timestamp)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Заголовки CSV
	headers := []string{"timestamp", "agent_id", "event_type", "severity", "user", "process", "message", "raw_log"}
	if err := writer.Write(headers); err != nil {
		return
	}

	// Данные
	for _, event := range events {
		row := []string{
			getString(event, "timestamp"),
			getString(event, "agent_id"),
			getString(event, "event_type"),
			getString(event, "severity"),
			getString(event, "user"),
			getString(event, "process"),
			getString(event, "message"),
			getString(event, "raw_log"),
		}
		if err := writer.Write(row); err != nil {
			return
		}
	}
}

func getString(m map[string]any, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}
