package main

import (
	"fmt"
	"log"

	"github.com/Narotan/Web-SIEM/Web/backend/internal/config"
	"github.com/Narotan/Web-SIEM/Web/backend/internal/repository"
	"github.com/Narotan/Web-SIEM/Web/backend/internal/service"
	transportHttp "github.com/Narotan/Web-SIEM/Web/backend/internal/transport/http"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()

	repo := repository.NewNosqlRepository(cfg.DBAddr)

	svc := service.NewSiemService(repo, cfg.DBName)

	handler := transportHttp.NewHandler(svc)

	r := gin.Default()

	transportHttp.SetupRouter(r, handler, cfg.WebUser, cfg.WebPass)

	r.Static("/css", "./frontend/css")
	r.Static("/js", "./frontend/js")

	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/index.html", "./frontend/index.html")
	r.StaticFile("/login.html", "./frontend/login.html")
	r.StaticFile("/events.html", "./frontend/events.html")

	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	log.Printf("Сервер запущен на http://localhost:%s", cfg.ServerPort)
	log.Printf("Frontend находится в ./frontend/")
	log.Printf("Логин: %s, Пароль: %s", cfg.WebUser, cfg.WebPass)

	if err := r.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
