package models

import (
	"time"
	"github.com/gofrs/uuid"
)

type Event struct {
	ID 			uuid.UUID 	`json:"id,omitempty"`
	AgendaIDs   []uuid.UUID `json:"agendaIds"`            // tableau d’UUIDs
	UID         string      `json:"uid"`
	Description string      `json:"description"`
	Name        string      `json:"name"` 				  //Summary
	Start       time.Time   `json:"start"`                // RFC3339 avec offset (ex: +01:00)
	End         time.Time   `json:"end"`
	Location    string      `json:"location"`
	LastUpdate  *time.Time  `json:"lastUpdate,omitempty"` // nullable dans l’iCal
}

type SchedulerPayload struct {
	AgendaID uuid.UUID 	`json:"agendaId"`
	Event    Event     	`json:"event"`
}

type FieldChange struct {
    Field  string `json:"field"`           		
    Before string `json:"before,omitempty"`	    
    After  string `json:"after,omitempty"`
}

type TimetableAlert struct {
    AgendaID uuid.UUID    `json:"agendaId"`
    UID      string       `json:"uid"`
    Changes  []FieldChange `json:"changes"`
    After    Event        `json:"after"`
    Before   *Event       `json:"before,omitempty"`
}
