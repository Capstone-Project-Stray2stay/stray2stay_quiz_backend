package adapter

import (
	"context"
	"errors"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/S-nudhana/stray2stay-quiz-backend/internal/core/domain"
)

type MongoQuizAdapter struct {
	db *mongo.Database
}

func NewMongoQuizAdapter(db *mongo.Database) *MongoQuizAdapter {
	return &MongoQuizAdapter{db: db}
}

func collectionForPetType(petType string) (string, error) {
	switch petType {
	case "cat":
		return "catBreeds", nil
	case "dog":
		return "dogBreeds", nil
	default:
		return "", domain.ErrInvalidPetType
	}
}

// RandomBreedWithImage picks one random breed document (via $sample) from the
// catBreeds/dogBreeds collection for petType, then picks a random image from
// that breed's possibleColors.
func (m *MongoQuizAdapter) RandomBreedWithImage(petType string) (string, string, error) {
	collectionName, err := collectionForPetType(petType)
	if err != nil {
		return "", "", err
	}

	collection := m.db.Collection(collectionName)

	// Only consider breeds that actually have at least one image, so we
	// never sample a document we'd have to throw away downstream.
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"possibleColors.0": bson.M{"$exists": true},
		}}},
		{{Key: "$sample", Value: bson.M{"size": 1}}},
	}

	ctx := context.Background()
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return "", "", err
	}
	defer cursor.Close(ctx)

	var docs []domain.BreedDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return "", "", err
	}

	if len(docs) == 0 {
		return "", "", domain.ErrNoBreedsFound
	}

	breed := docs[0]
	if len(breed.PossibleColors) == 0 {
		return "", "", domain.ErrNoImageForBreed
	}

	pick := breed.PossibleColors[rand.Intn(len(breed.PossibleColors))]
	if pick.Image == "" {
		return "", "", errors.New("sampled image is empty")
	}

	return breed.BreedName, pick.Image, nil
}

// RandomDistractorBreedNames samples up to count breed names from the
// catBreeds/dogBreeds collection for petType, excluding excludeBreed, to use
// as wrong multiple-choice options alongside the correct breed.
func (m *MongoQuizAdapter) RandomDistractorBreedNames(petType string, excludeBreed string, count int) ([]string, error) {
	collectionName, err := collectionForPetType(petType)
	if err != nil {
		return nil, err
	}

	collection := m.db.Collection(collectionName)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"breedName": bson.M{"$ne": excludeBreed}}}},
		{{Key: "$sample", Value: bson.M{"size": count}}},
		{{Key: "$project", Value: bson.M{"breedName": 1, "_id": 0}}},
	}

	ctx := context.Background()
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []domain.BreedDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(docs))
	for _, doc := range docs {
		names = append(names, doc.BreedName)
	}

	return names, nil
}
