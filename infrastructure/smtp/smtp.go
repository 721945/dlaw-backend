package smtp

import (
	"fmt"
	"github.com/721945/dlaw-backend/libs"
	"github.com/go-gomail/gomail"
	"google.golang.org/api/calendar/v3"
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

func (c SMTP) SendCalendarInvitation(email string, event *calendar.Event) error {
	// create a new email message
	m := gomail.NewMessage()

	// set the subject and body of the email message
	m.SetHeader("From", c.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", event.Summary)

	// create a new multipart message
	m.SetBody("plain", event.HtmlLink)

	// send the email message
	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Sent email to " + email)

	return nil
}

//
//func (c SMTP) SendCalendarInvitationLinks(email string, event *calendar.Event) error {
//	// create a new email message
//	m := gomail.NewMessage()
//
//	// set the subject and body of the email message
//	m.SetHeader("From", c.Username)
//	m.SetHeader("To", email)
//	m.SetHeader("Subject", "Your OTP code")
//	m.SetHeader("Subject", event.Summary)
//	// create a new buffer to hold the iCalendar data
//	buf := new(bytes.Buffer)
//
//	// encode the iCalendar data to the buffer
//	if err := event.Encode(buf); err != nil {
//		return fmt.Errorf("unable to encode iCalendar data: %v", err)
//	}
//
//	// attach the iCalendar data to the email message as a MIME part
//	m.SetBody("text/calendar; charset=UTF-8; method=REQUEST", buf.String())
//
//	//m.SetBody()
//
//	// create a new multipart message
//	//mp := multipart.New("mixed")
//	//
//	//// create a new part for the iCalendar attachment
//	//calPart := multipart.NewPart(nil)
//	//calPart.Header.Set("Content-Type", "text/calendar; charset=UTF-8; method=REQUEST")
//	//calPart.Header.Set("Content-Transfer-Encoding", "base64")
//	//calPart.SetFileName("event.ics")
//
//	// encode the iCalendar data and write it to the part
//	//enc := base64.NewEncoder(base64.StdEncoding, calPart)
//	//if err := event.Encode(enc); err != nil {
//	//	return fmt.Errorf("unable to encode iCalendar attachment: %v", err)
//	//}
//	//enc.Close()
//	//
//	//// add the part to the multipart message
//	//mp.AddPart(calPart)
//	//
//	//// attach the multipart message to the email message
//	//buf := new(bytes.Buffer)
//	//mp.WriteTo(buf)
//	//m.SetBody("multipart/mixed", buf.String())
//
//	// send the email message
//	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
//	if err := d.DialAndSend(m); err != nil {
//		return err
//	}
//
//	return nil
//}
