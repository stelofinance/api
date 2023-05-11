package tools

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type envVariables struct {
	DbConnection  string
	ProductionEnv bool
	AdminKey      string
	Pusher        PusherVariables
}

type PusherVariables struct {
	Host      string
	AppId     string
	AppKey    string
	AppSecret string
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
	adminKey, err := getEnvVariable("ADMIN_API_KEY")
	if err != nil {
		return err
	}
	pusherHost, err := getEnvVariable("PUSHER_HOST")
	if err != nil {
		return err
	}
	pusherAppId, err := getEnvVariable("PUSHER_APP_ID")
	if err != nil {
		return err
	}
	pusherAppKey, err := getEnvVariable("PUSHER_APP_KEY")
	if err != nil {
		return err
	}
	pusherAppSecret, err := getEnvVariable("PUSHER_APP_SECRET")
	if err != nil {
		return err
	}
	EnvVars = envVariables{
		DbConnection:  dbConnectionString,
		ProductionEnv: prodEnv,
		AdminKey:      adminKey,
		Pusher: PusherVariables{
			Host:      pusherHost,
			AppId:     pusherAppId,
			AppKey:    pusherAppKey,
			AppSecret: pusherAppSecret,
		},
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