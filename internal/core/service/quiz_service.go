package service

import (
	"strings"
	"time"

	"github.com/S-nudhana/stray2stay-quiz-backend/internal/core/domain"
	"github.com/S-nudhana/stray2stay-quiz-backend/internal/core/port"
)

type quizService struct {
	repo   port.QuizMongoRepository
	signer *tokenSigner
	ttl    time.Duration
}

func NewQuizService(repo port.QuizMongoRepository, tokenSecret string, ttlMinutes int) port.QuizService {
	return &quizService{
		repo:   repo,
		signer: newTokenSigner(tokenSecret),
		ttl:    time.Duration(ttlMinutes) * time.Minute,
	}
}

func (s *quizService) NewQuestion(petType string) (string, domain.QuizQuestion, error) {
	normalized := strings.ToLower(strings.TrimSpace(petType))
	if normalized != "cat" && normalized != "dog" {
		return "", domain.QuizQuestion{}, domain.ErrInvalidPetType
	}

	breedName, imageURL, err := s.repo.RandomBreedWithImage(normalized)
	if err != nil {
		return "", domain.QuizQuestion{}, err
	}

	now := time.Now().Unix()
	token, err := s.signer.sign(tokenPayload{
		PetType:   normalized,
		BreedName: breedName,
		ImageURL:  imageURL,
		IssuedAt:  now,
	})
	if err != nil {
		return "", domain.QuizQuestion{}, err
	}

	question := domain.QuizQuestion{
		PetType:  normalized,
		ImageURL: imageURL,
		IssuedAt: now,
	}

	return token, question, nil
}

func (s *quizService) CheckAnswer(token string, guess string) (domain.QuizAnswerResult, error) {
	payload, err := s.signer.verify(token)
	if err != nil {
		return domain.QuizAnswerResult{}, err
	}

	if time.Since(time.Unix(payload.IssuedAt, 0)) > s.ttl {
		return domain.QuizAnswerResult{}, domain.ErrInvalidToken
	}

	correct := isCloseEnough(guess, payload.BreedName)

	return domain.QuizAnswerResult{
		Correct:      correct,
		CorrectBreed: payload.BreedName,
	}, nil
}
