package services

import (
	"errors"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/golang-jwt/jwt/v4"
)

type JWTAuthService struct {
	env    libs.Env
	logger libs.Logger
}

func NewJWTAuthService(env libs.Env, logger libs.Logger) JWTAuthService {
	return JWTAuthService{
		env:    env,
		logger: logger,
	}
}

func (j JWTAuthService) GenerateToken(user models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.Firstname,
		"last_name":  user.Lastname,
		//
		//"exp":        time.Now().Add(time.Hour * 24).Unix(),

	})

	tokenString, err := token.SignedString(j.env.JWTSecret)

	if err != nil {
		j.logger.Error(err)
	}

	return tokenString
}

func (j JWTAuthService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.env.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
