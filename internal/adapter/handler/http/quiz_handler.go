package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/S-nudhana/stray2stay-quiz-backend/internal/core/domain"
	"github.com/S-nudhana/stray2stay-quiz-backend/internal/core/port"
)

type QuizHandler struct {
	service port.QuizService
}

func NewQuizHandler(service port.QuizService) *QuizHandler {
	return &QuizHandler{service: service}
}

type answerRequest struct {
	Token string `json:"token"`
	Guess string `json:"guess"`
}

func (h *QuizHandler) GetQuestion(c *fiber.Ctx) error {
	petType := c.Query("petType")

	token, question, err := h.service.NewQuestion(petType)
	if err != nil {
		return respondQuizError(c, err)
	}

	return c.JSON(fiber.Map{
		"token":    token,
		"petType":  question.PetType,
		"imageUrl": question.ImageURL,
		"message":  "Get quiz question successfully",
	})
}

func (h *QuizHandler) SubmitAnswer(c *fiber.Ctx) error {
	payload := new(answerRequest)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	if payload.Token == "" || payload.Guess == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "token and guess are required"})
	}

	result, err := h.service.CheckAnswer(payload.Token, payload.Guess)
	if err != nil {
		return respondQuizError(c, err)
	}

	return c.JSON(fiber.Map{
		"correct":      result.Correct,
		"correctBreed": result.CorrectBreed,
		"message":      "Checked answer successfully",
	})
}

func respondQuizError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidPetType):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, domain.ErrInvalidToken):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, domain.ErrNoBreedsFound), errors.Is(err, domain.ErrNoImageForBreed):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Something went wrong"})
	}
}
