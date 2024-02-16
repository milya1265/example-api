package service

import (
	repository "example1/internal/repository/sqlc/generate"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(password []byte) ([]byte, error) {
	hashpassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashpassword, nil
}

func ComparePassword(password, hashpassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashpassword, password)
	return err
}

func NewToken(user repository.User, role int32, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["login"] = user.Login
	claims["role"] = role
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseSubject(Token string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(Token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, TakeClaimsErr
	}

	return claims, nil
}

func ParseExpiration(t string, s string) (int64, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, TakeClaimsErr
	}

	return int64(claims["exp"].(float64)), nil
}
