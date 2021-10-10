package objects

import (
	"encoding/json"
	"net/http"
)

// MaxListLimit maximum listing
const MaxListLimit = 200

// GetRequest for retreiving single Event
type GetRequest struct {
	ID string `json:"id"`
}

// ListRequest for retreiving list of Events
type ListRequest struct {
	Limit int    `json:"limit"`
	After string `json:"after"`
	// Optional name matching
	Name string `json:"name"`
}

// CreateRequest for creating a new Event
type CreateRequest struct {
	Event *Event `json:"event"`
}

// UpdateDetailsRequest to update existing Event
type UpdateDetailsRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

// CancelRequest to cancel an Event
type CancelRequest struct {
	ID string `json:"id"`
}

// RescheduleRequest to reschedule an Event
type RescheduleRequest struct {
	ID      string    `json:"id"`
	NewSlot *TimeSlot `json:"new_slot"`
}

// DeleteRequest to delete an Event
type DeleteRequest struct {
	ID string `json:"id"`
}

// EventResponseWrapper response od any Event request
type EventResponseWrapper struct {
	Event  *Event   `json:"event"`
	Events []*Event `json:"events,omitempty"`
	Code   int      `json:"-"`
}

func (e *EventResponseWrapper) Json() []byte {
	if e == nil {
		return []byte("{}")
	}
	res, _ := json.Marshal(e)
	return res
}

func (e *EventResponseWrapper) StatusCode() int {
	if e == nil || e.Code == 0 {
		return http.StatusOK
	}
	return e.Code
}
