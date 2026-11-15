package domain

import "errors"

var (
	ErrInvalidPetType  = errors.New("petType must be \"cat\" or \"dog\"")
	ErrNoBreedsFound   = errors.New("no breeds found for pet type")
	ErrNoImageForBreed = errors.New("breed has no image")
	ErrInvalidToken    = errors.New("invalid or expired question token")
)

// BreedColor mirrors a single entry of a breed document's possibleColors
// array in the catBreeds/dogBreeds Mongo collections.
type BreedColor struct {
	Color string `bson:"color"`
	Image string `bson:"link"`
}

// BreedDoc mirrors the shape of documents in catBreeds/dogBreeds.
type BreedDoc struct {
	BreedName      string       `bson:"breedName"`
	PossibleColors []BreedColor `bson:"possibleColors"`
}

// QuizQuestion is the payload encoded into the signed question token and
// returned (minus the answer) to the client.
type QuizQuestion struct {
	PetType   string `json:"petType"`
	BreedName string `json:"-"`
	ImageURL  string `json:"imageUrl"`
	IssuedAt  int64  `json:"-"`
}

type QuizAnswerResult struct {
	Correct      bool   `json:"correct"`
	CorrectBreed string `json:"correctBreed"`
}
