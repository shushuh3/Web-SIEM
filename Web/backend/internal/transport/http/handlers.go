package http

import (
	"net/http"
	"strconv"

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
