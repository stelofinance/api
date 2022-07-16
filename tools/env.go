package tools

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}

func GetEnvVariable(valueToGet string) (string, error) {
	value := os.Getenv(valueToGet)
	if value == "" {
		return "", errors.New("Environment variable not found")
	}
	return value, nil
}