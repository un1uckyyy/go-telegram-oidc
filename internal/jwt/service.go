package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/un1uckyyy/go-telegram-oidc/pkg/logger"
)

func GetSubjectFromJwt(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		logger.ErrorLogger.Printf("error parsing jwt token: %v", err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if sub, exists := claims["sub"]; exists {
			return sub.(string), nil
		}

		return "", fmt.Errorf("error getting subject from jwt token")
	}

	return "", fmt.Errorf("jwt token is not valid")
}
