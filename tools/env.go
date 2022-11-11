package tools

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type envVariables struct {
	DbConnection  string
	ProductionEnv bool
	JwtSecret     []byte
	AdminKey      string
}

var EnvVars envVariables

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	dbConnectionString, err := getEnvVariable("DB_CONNECTION_STRING")
	if err != nil {
		return err
	}
	prodEnvString, err := getEnvVariable("PRODUCTION_ENV")
	var prodEnv bool
	if err != nil {
		return err
	} else {
		prodEnv, err = strconv.ParseBool(prodEnvString)
		if err != nil {
			return err
		}
	}
	jwtSecret, err := getEnvVariable("JWT_SECRET")
	if err != nil {
		return err
	}
	adminKey, err := getEnvVariable("ADMIN_KEY")
	if err != nil {
		return err
	}
	EnvVars = envVariables{
		DbConnection:  dbConnectionString,
		ProductionEnv: prodEnv,
		JwtSecret:     []byte(jwtSecret),
		AdminKey:      adminKey,
	}
	return nil
}

func getEnvVariable(valueToGet string) (string, error) {
	value := os.Getenv(valueToGet)
	if value == "" {
		return "", errors.New("Environment variable not found")
	}
	return value, nil
}