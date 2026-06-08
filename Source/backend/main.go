package main

import (
	"fmt"
	"harbor/config"
	"harbor/handlers"
	"harbor/repository"
	"harbor/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	var repo repository.Repository
	if cfg.UseDB {
		log.Println("Хранилище: PostgreSQL")
		pg, err := repository.NewPostgresRepository(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		if err != nil {
			log.Fatalf("Ошибка подключения к БД: %v", err)
		}
		defer pg.Close()
		repo = pg
	} else {
		log.Println("Хранилище: In-Memory (данные не сохраняются)")
		repo = repository.NewMemoryRepository()
	}

	svc := service.New(repo, cfg.JWTSecret, cfg.JWTExpiryHours)
	h := handlers.New(svc)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	h.Register(r)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Сервер: http://localhost%s", addr)
	log.Println("Аккаунты: admin@harbor.ru/admin123 | dispatcher@harbor.ru/disp123 | operator@harbor.ru/oper123")
	r.Run(addr)
}
