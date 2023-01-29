package tools

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type envVariables struct {
	DbConnection     string
	ProductionEnv    bool
	JwtSecret        []byte
	AdminKey         string
	CentrifugoApiKey string
	CentrifugoAddr   string
	CentrifugoJWTKey []byte
}

var EnvVars envVariables

func isProduction() bool {
	prodEnvVar := os.Getenv("PRODUCTION_ENV")
	switch prodEnvVar {
	case "1", "t", "T", "true", "TRUE", "True":
		return true
	}
	return false
}

func LoadEnv() error {
	prodEnv := isProduction()
	if !prodEnv {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}
	dbConnectionString, err := getEnvVariable("DB_CONNECTION_STRING")
	if err != nil {
		return err
	}
	jwtSecret, err := getEnvVariable("JWT_SECRET")
	if err != nil {
		return err
	}
	adminKey, err := getEnvVariable("ADMIN_API_KEY")
	if err != nil {
		return err
	}
	centrifugoKey, err := getEnvVariable("CENTRIFUGO_API_KEY")
	if err != nil {
		return err
	}
	centrifugoAddr, err := getEnvVariable("CENTRIFUGO_API_ADDR")
	if err != nil {
		return err
	}
	centrifugoJwtKey, err := getEnvVariable("CENTRIFUGO_JWT_KEY")
	if err != nil {
		return err
	}
	EnvVars = envVariables{
		DbConnection:     dbConnectionString,
		ProductionEnv:    prodEnv,
		JwtSecret:        []byte(jwtSecret),
		AdminKey:         adminKey,
		CentrifugoApiKey: centrifugoKey,
		CentrifugoAddr:   centrifugoAddr,
		CentrifugoJWTKey: []byte(centrifugoJwtKey),
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