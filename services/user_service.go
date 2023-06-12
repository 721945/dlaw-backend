package services

import (
	"fmt"
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/api/utils"
	"github.com/721945/dlaw-backend/infrastructure/smtp"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type UserService interface {
	WithTrx(trxHandle *gorm.DB) UserService
	CreateUser(user models.User) (models.User, error)
	UpdateUser(id uuid.UUID, user models.User) error
	DeleteUser(id uuid.UUID) error
	GetUser(id uuid.UUID) (*dtos.UserDto, error)
	GetUsers(organization string) ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByEmailAndOrganization(email, organization string) (*models.User, error)
	ForgetPassword(email, organization string) error
	VerifyOTP(email, organization string, otp string) error
	ResetPassword(email, organization, otp, newPassword string) error
	ChangedPassword(id uuid.UUID, current string, newPassword string) error
}

type userService struct {
	userRepo repositories.UserRepository
	logger   *libs.Logger
	smtp     smtp.SMTP
}

func NewUserService(userRepo repositories.UserRepository, logger *libs.Logger, smtp smtp.SMTP) UserService {
	return &userService{userRepo: userRepo, logger: logger, smtp: smtp}
}

func (u userService) WithTrx(trxHandle *gorm.DB) UserService {
	if trxHandle == nil {
		return u
	}

	u.userRepo = u.userRepo.WithTrx(trxHandle)
	return u
}

func (u userService) CreateUser(user models.User) (models.User, error) {
	return u.userRepo.CreateUser(user)
}

func (u userService) UpdateUser(id uuid.UUID, user models.User) error {
	return u.userRepo.UpdateUser(id, user)
}

func (u userService) DeleteUser(id uuid.UUID) error {
	return u.userRepo.DeleteUser(id)
}

func (u userService) GetUser(id uuid.UUID) (*dtos.UserDto, error) {
	user, err := u.userRepo.GetUser(id)

	if err != nil {
		return nil, err
	}

	return &dtos.UserDto{
		ID:        user.ID.String(),
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
	}, nil
}

func (u userService) GetUsers(organization string) ([]models.User, error) {
	if organization != "" {
		return u.userRepo.GetUsersByOrganization(organization)
	}
	return u.userRepo.GetUsers()
}

func (u userService) GetUserByEmail(email string) (user *models.User, err error) {
	return u.userRepo.GetUserByEmail(email)
}
func (u userService) GetUserByEmailAndOrganization(email, organization string) (user *models.User, err error) {
	return u.userRepo.GetUserByEmail(email)
}

func (u userService) ForgetPassword(email, org string) error {
	user, err := u.userRepo.GetUserByEmailAndOrganization(email, org)

	if err != nil {
		return err
	}

	otp := generateOTP()
	user.OtpSecret = &otp
	timeLeft := time.Now().Add(time.Minute * 30)
	user.OtpExpiredAt = &timeLeft

	err = u.userRepo.UpdateUser(user.ID, *user)

	if err != nil {
		return err
	}

	err = u.smtp.SendOTPtoEmail(email, otp, 5)

	if err != nil {
		return err
	}

	// send email
	return nil
}

func (u userService) ChangedPassword(id uuid.UUID, current string, new string) error {
	user, err := u.userRepo.GetUser(id)

	if err != nil {
		return err
	}

	isSame := utils.CheckPasswordHash(current, user.Password)

	if !isSame {
		return fmt.Errorf("wrong password")
	}

	hashedPassword, err := utils.HashPassword(new)

	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return u.userRepo.UpdateUser(id, *user)

}

func (u userService) VerifyOTP(email, organization, otp string) error {
	_, err := u.verifyOtp(email, organization, otp)
	return err
}

func (u userService) ResetPassword(email, organization, otp, newPassword string) error {
	user, err := u.verifyOtp(email, organization, otp)

	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)

	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return u.userRepo.UpdateUser(user.ID, *user)
}

func (u userService) verifyOtp(email, organization, otp string) (*models.User, error) {
	user, err := u.userRepo.GetUserByEmailAndOrganization(email, organization)

	if err != nil {
		return nil, err
	}

	if user.OtpSecret == nil {
		return nil, fmt.Errorf("otp not found")
	}

	if *user.OtpSecret != otp {
		return nil, fmt.Errorf("wrong otp")
	}

	if user.OtpExpiredAt.Before(time.Now()) {
		return nil, fmt.Errorf("otp expired")
	}

	return user, err
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999)
	return fmt.Sprintf("%06d", otp)
}
