package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	PostgresConnectionString string
	JWTSecretKey             string
	JWTExpirationInSeconds   int64
	RequestMethod            string
}

var Envs = initConfig()

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

func initConfig() Config {
	_ = godotenv.Load()
	return Config{
		PostgresConnectionString: getEnv("POSTGRES_CONNECTION_STRING", "postgres"),
		JWTSecretKey:             getEnv("JWT_SECRET", "verySecureRandomKey"),
		JWTExpirationInSeconds:   getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		RequestMethod:            getEnv("REQUEST_METHOD", "Bearer"),
	}
}
