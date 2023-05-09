package google_calendar

import (
	"context"
	"fmt"
	"github.com/721945/dlaw-backend/infrastructure/smtp"
	"github.com/721945/dlaw-backend/libs"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
)

type GoogleCalendar struct {
	srv  *calendar.Service
	smtp smtp.SMTP
}

func NewGoogleCalendar(env libs.Env, smtp smtp.SMTP) GoogleCalendar {
	srv, err := getCalendarClient(env.GoogleCredPath)

	if err != nil {
		panic(err)
	}

	return GoogleCalendar{
		srv:  srv,
		smtp: smtp,
	}
}

func getCalendarClient(path string) (*calendar.Service, error) {
	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithCredentialsFile(path))
	if err != nil {
		return nil, fmt.Errorf("unable to create calendar service: %v", err)
	}
	return srv, nil
}

//func (c GoogleCalendar) CreateEvent(title, startTime, endTime string, emails []string, location, description *string) (*calendar.Event, error) {
//	attendees := make([]*calendar.EventAttendee, len(emails))
//
//	for i, email := range emails {
//		attendees[i] = &calendar.EventAttendee{Email: email, ResponseStatus: "needsAction"}
//	}
//
//	event := &calendar.Event{
//		Summary:     title,
//		Location:    *location,
//		Description: *description,
//		Start:       &calendar.EventDateTime{DateTime: startTime, TimeZone: "Asia/Bangkok"},
//		End:         &calendar.EventDateTime{DateTime: endTime, TimeZone: "Asia/Bangkok"},
//		Attendees:   attendees,
//	}
//
//	createdEvent, err := c.srv.Events.Insert("primary", event).Do()
//	if err != nil {
//		return nil, fmt.Errorf("unable to create event: %v", err)
//	}
//
//	return createdEvent, nil
//}

func (c GoogleCalendar) CreateEvent(title, startTime, endTime string, emails []string, location, description *string) (*calendar.Event, error) {
	attendees := make([]*calendar.EventAttendee, len(emails))
	for i, email := range emails {
		attendees[i] = &calendar.EventAttendee{
			Email:          email,
			ResponseStatus: "needsAction",
		}
	}

	// Set the calendar ID to the email address of the service account
	calendarID := "dlaw-service-2@dlaw-dev.iam.gserviceaccount.com"

	event := &calendar.Event{
		Summary:     title,
		Location:    *location,
		Description: *description,
		Start:       &calendar.EventDateTime{DateTime: startTime, TimeZone: "Asia/Bangkok"},
		End:         &calendar.EventDateTime{DateTime: endTime, TimeZone: "Asia/Bangkok"},
		Attendees:   attendees,
	}
	// Insert the event into the service account's calendar
	createdEvent, err := c.srv.Events.Insert(calendarID, event).SendUpdates("all").Do()
	if err != nil {
		log.Fatalf("Unable to create event: %v", err)
	}

	fmt.Printf("Event created with ID: %s\n", createdEvent.Id)

	return createdEvent, nil
}

func (c GoogleCalendar) UpdateEvent(event *calendar.Event) (*calendar.Event, error) {
	updatedEvent, err := c.srv.Events.Update("primary", event.Id, event).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to update event: %v", err)
	}

	return updatedEvent, nil
}

func (c GoogleCalendar) DeleteEvent(eventID string) error {
	err := c.srv.Events.Delete("primary", eventID).Do()
	if err != nil {
		return fmt.Errorf("unable to delete event: %v", err)
	}

	return nil
}
