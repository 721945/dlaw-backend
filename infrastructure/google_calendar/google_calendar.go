package google_calendar

import (
	"context"
	"fmt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendar struct {
	srv *calendar.Service
}

func NewGoogleCalendar() GoogleCalendar {
	srv, err := getCalendarClient("/Users/iam721945/.config/gcloud/application_default_credentials.json")

	if err != nil {
		panic(err)
	}

	return GoogleCalendar{
		srv: srv,
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

func (c GoogleCalendar) CreateEvent(title, startTime, endTime string, emails []string, location, description *string) (*calendar.Event, error) {
	attendees := make([]*calendar.EventAttendee, len(emails))

	for i, email := range emails {
		attendees[i] = &calendar.EventAttendee{Email: email}
	}

	event := &calendar.Event{
		Summary:     title,
		Location:    *location,
		Description: *description,
		Start:       &calendar.EventDateTime{DateTime: startTime, TimeZone: "Asia/Bangkok"},
		End:         &calendar.EventDateTime{DateTime: endTime, TimeZone: "Asia/Bangkok"},
		Attendees:   attendees,
	}

	createdEvent, err := c.srv.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create event: %v", err)
	}

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
