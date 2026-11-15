package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	adapter "github.com/S-nudhana/stray2stay-quiz-backend/internal/adapter/database"
	handler "github.com/S-nudhana/stray2stay-quiz-backend/internal/adapter/handler/http"
	"github.com/S-nudhana/stray2stay-quiz-backend/internal/adapter/handler/router"
	"github.com/S-nudhana/stray2stay-quiz-backend/internal/core/service"
	"github.com/S-nudhana/stray2stay-quiz-backend/internal/infrastructure/config"
	"github.com/S-nudhana/stray2stay-quiz-backend/internal/infrastructure/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mongoClient, err := database.NewMongoDatabase(cfg.Mongo)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	mongoDB := mongoClient.Database(cfg.Mongo.DBName)

	quizRepo := adapter.NewMongoQuizAdapter(mongoDB)
	quizService := service.NewQuizService(quizRepo, cfg.Quiz.TokenSecret, cfg.Quiz.TokenTTLMins)
	quizHandler := handler.NewQuizHandler(quizService)

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORS.AllowOrigins,
		AllowMethods: cfg.CORS.AllowMethods,
		AllowHeaders: cfg.CORS.AllowHeaders,
		MaxAge:       cfg.CORS.MaxAge,
	}))

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	router.RegisterQuizRoutes(app, quizHandler)

	log.Printf("stray2stay-quiz-backend listening on %s", cfg.Server.Addr)
	if err := app.Listen(cfg.Server.Addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
