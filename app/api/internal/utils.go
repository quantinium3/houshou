package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"strconv"
	"time"
)

type JWTClaims struct {
	UserId string `json:"name"`
	Role   string `json:"role"`
	Admin  bool   `json:"admin"`
	jwt.RegisteredClaims
}

func CreateNamespacedId(name string) string {
	return fmt.Sprintf("%v_%v", name, uuid.New().String())
}

func GetEnvOr(env string, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		return defaultValue
	}
	return value
}

func CreateAccessToken(id string) (string, error) {
	accessExpiry, err := strconv.ParseInt(GetEnvOr("ACCESS_EXPIRY", "60"),
		10, 64)
	if err != nil {
		return "", err
	}

	accessSecret := os.Getenv("ACCESS_SECRET")
	if accessSecret == "" {
		return "", errors.New("ACCESS_SECRET environment variable not set")
	}

	claims := JWTClaims{
		UserId: id,
		Role:   "user",
		Admin:  false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.
				Hour * time.Duration(accessExpiry))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func CreateRefreshToken(id string) (string, error) {
	refreshExpiry, err := strconv.ParseInt(GetEnvOr("REFRESH_EXPIRY", "30"),
		10, 64)
	if err != nil {
		return "", err
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		return "", errors.New("REFRESH_SECRET environment variable not set")
	}

	claims := JWTClaims{
		UserId: id,
		Role:   "user",
		Admin:  false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * time.
				Duration(refreshExpiry))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
