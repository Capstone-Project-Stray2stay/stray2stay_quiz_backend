package router

import (
	"github.com/gofiber/fiber/v2"

	handler "github.com/S-nudhana/stray2stay-quiz-backend/internal/adapter/handler/http"
)

func RegisterQuizRoutes(app *fiber.App, quizHandler *handler.QuizHandler) {
	quiz := app.Group("/api/quiz")

	quiz.Get("/question", quizHandler.GetQuestion)
	quiz.Post("/answer", quizHandler.SubmitAnswer)
}
