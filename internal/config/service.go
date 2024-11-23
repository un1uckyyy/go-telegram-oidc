package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

var Instance *Config

type Config struct {
	TgToken      string `validate:"required"`
	ClientId     string `validate:"required"`
	ClientSecret string `validate:"required"`
	RedirectHost string `validate:"required"`
	AuthUrl      string `validate:"required"`
	TokenUrl     string `validate:"required"`
}

func LoadConfig() error {
	tgToken := os.Getenv("TG_TOKEN")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectHost := os.Getenv("REDIRECT_HOST")
	authUrl := os.Getenv("AUTH_URL")
	tokenUrl := os.Getenv("TOKEN_URL")

	config := &Config{
		TgToken:      tgToken,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectHost: redirectHost,
		AuthUrl:      authUrl,
		TokenUrl:     tokenUrl,
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	Instance = config

	return nil
}
