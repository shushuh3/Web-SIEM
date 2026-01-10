package http

import (
	"github.com/Narotan/Web-SIEM/Web/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, handler *Handler, webUser, webPass string) {
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/health", middleware.BasicAuth(webUser, webPass), handler.Health)

		protected := api.Group("")
		protected.Use(middleware.BasicAuth(webUser, webPass))
		{
			protected.GET("/events", handler.GetEvents)
			protected.GET("/events/export", handler.ExportEvents)
			protected.GET("/stats", handler.GetStats)
		}
	}
}
