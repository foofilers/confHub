package auth

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/spf13/viper"
	"errors"
	"time"
	"github.com/sirupsen/logrus"
)

func ValidateJwt(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("Invalid token")
	}
	if standardClaims, ok := token.Claims.(*jwt.StandardClaims); ok {
		return standardClaims.Subject, nil
	}
	return "", errors.New("Cannot parse JWT claims")
}

func GenerateJwt(userId string, duration time.Duration) (string, error) {
	stdClaims := &jwt.StandardClaims{
		Issuer:"cfh",
		Subject:userId,
		IssuedAt:time.Now().Unix(),
	}
	if duration > 0 {
		stdClaims.ExpiresAt = time.Now().Add(duration).Unix()
	}
	logrus.Debugf("generated new jwtToken claims:%+v", stdClaims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
	return token.SignedString([]byte(viper.GetString("jwt.secret")))
}