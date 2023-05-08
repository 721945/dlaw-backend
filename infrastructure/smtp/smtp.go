package smtp

import (
	"fmt"
	"github.com/721945/dlaw-backend/libs"
	"github.com/go-gomail/gomail"
)

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTP(env libs.Env) SMTP {
	return SMTP{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: env.SMTPAddress,
		Password: env.SMTPPassword,
	}
}

func (c SMTP) SendOTPtoEmail(email string, otp string, minute int) error {
	// create a new email message
	m := gomail.NewMessage()

	// set the subject and body of the email message
	m.SetHeader("From", c.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your OTP code")
	m.SetBody("text/plain", "Your OTP code is "+otp+" and it will be expired in "+fmt.Sprint(minute)+" minutes")

	// create a new SMTP dialer
	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)

	// send the email message
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
