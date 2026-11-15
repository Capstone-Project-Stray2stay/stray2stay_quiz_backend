package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	CORS   CORSConfig
	Mongo  MongoConfig
	Quiz   QuizConfig
}

type ServerConfig struct {
	Env  string
	Addr string
}

func (s ServerConfig) IsProduction() bool {
	return s.Env == "production"
}

type CORSConfig struct {
	AllowOrigins string
	AllowMethods string
	AllowHeaders string
	MaxAge       int
}

type MongoConfig struct {
	URI    string
	DBName string
}

// QuizConfig configures the signed question tokens that let the backend stay
// stateless between GET /question and POST /answer (no session store, so a
// fresh server instance can still verify a token issued before it restarted).
type QuizConfig struct {
	TokenSecret  string
	TokenTTLMins int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	v := &validator{}

	env := getEnv("ENV", "development")

	cfg := &Config{
		Server: ServerConfig{
			Env:  env,
			Addr: getEnv("ADDR", ":8081"),
		},
		CORS: CORSConfig{
			AllowOrigins: getEnv("ORIGIN", "*"),
			AllowMethods: "GET,POST,OPTIONS",
			AllowHeaders: "Origin,Content-Type,Accept",
			MaxAge:       300,
		},
		Mongo: MongoConfig{
			URI:    v.require("MONGO_DB_URI"),
			DBName: v.require("MONGO_DB_NAME"),
		},
		Quiz: QuizConfig{
			TokenSecret:  v.require("QUIZ_TOKEN_SECRET"),
			TokenTTLMins: getEnvInt("QUIZ_TOKEN_TTL_MINUTES", 10),
		},
	}

	if err := v.err(); err != nil {
		return nil, err
	}

	return cfg, nil
}

type validator struct {
	missing []string
}

func (v *validator) require(key string) string {
	val := os.Getenv(key)
	if val == "" {
		v.missing = append(v.missing, key)
	}
	return val
}

func (v *validator) err() error {
	if len(v.missing) == 0 {
		return nil
	}
	return fmt.Errorf("missing required environment variables: %v", v.missing)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	var i int
	if _, err := fmt.Sscanf(val, "%d", &i); err != nil {
		return fallback
	}
	return i
}
