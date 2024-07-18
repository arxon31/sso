package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	envSecretKey = "SECRET_KEY"
	duration     = 30 * time.Minute
)

var secret = secretKeyFromEnv()

// NewToken creates new JWT token for given username
func NewToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func secretKeyFromEnv() string {
	key, exists := os.LookupEnv(envSecretKey)
	if !exists {
		panic("secret key is not provided")
	}

	return key
}
