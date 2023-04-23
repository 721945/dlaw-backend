package dtos

import "github.com/721945/dlaw-backend/models"

type CreateUserDto struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type UpdateUserDto struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type UserDto struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func ToUserDto(user *models.User) *UserDto {
	return &UserDto{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.Firstname,
		LastName:  user.Lastname,
	}
}
