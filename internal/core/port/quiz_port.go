package port

import "github.com/S-nudhana/stray2stay-quiz-backend/internal/core/domain"

// QuizMongoRepository reads breed reference data from the shared
// stray2stay MongoDB (catBreeds/dogBreeds collections).
type QuizMongoRepository interface {
	RandomBreedWithImage(petType string) (breedName string, imageURL string, err error)
}

type QuizService interface {
	NewQuestion(petType string) (token string, question domain.QuizQuestion, err error)
	CheckAnswer(token string, guess string) (domain.QuizAnswerResult, error)
}
