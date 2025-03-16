package auth

import (
	"github.com/Code-xCartel/noxus-api-svc/config"
	"github.com/Code-xCartel/noxus-api-svc/types/auth"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateJWT(user *auth.User) (string, error) {
	secret := config.Envs.JWTSecretKey
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"noxId":  user.NoxID,
		"exp":    time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
