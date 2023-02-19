package services

import (
	"errors"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type JWTAuthService struct {
	env    libs.Env
	logger *libs.Logger
}

func NewJWTAuthService(env libs.Env, logger *libs.Logger) JWTAuthService {
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
	})

	tokenString, err := token.SignedString([]byte(j.env.JWTSecret))

	if err != nil {
		j.logger.Error(err)
	}

	return tokenString
}

func (j JWTAuthService) Authorize(tokenString string) (bool, error) {
	_, err := j.VerifyToken(tokenString)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (j JWTAuthService) VerifyToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.env.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, err
	}

	return &models.User{
		Model:     gorm.Model{ID: uint(claims["id"].(float64))},
		Email:     claims["email"].(string),
		Firstname: claims["first_name"].(string),
		Lastname:  claims["last_name"].(string),
	}, nil
}

//

func (j JWTAuthService) GetUserIDFromToken(tokenString string) (uint, error) {
	user, err := j.VerifyToken(tokenString)

	if err != nil {
		return 0, err
	}

	//id := user["id"].(float64)
	id := user.ID

	return uint(id), nil
}
