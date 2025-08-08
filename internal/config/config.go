package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MongoURI  string
	DBName    string
	JWTSecret string
}

func Load() Config {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	db := os.Getenv("DB_NAME")
	if db == "" {
		db = "expense_db"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change_me"
	}
	return Config{Port: port, MongoURI: uri, DBName: db, JWTSecret: secret}
}
