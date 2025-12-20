package models

import (
	"time"
	"github.com/gofrs/uuid"
)

type Event struct {
	ID          uuid.UUID   `json:"id"`
	AgendaIDs   []uuid.UUID `json:"agendaIds"`            // tableau d’UUIDs
	UID         string      `json:"uid"`
	Description string      `json:"description"`
	Name        string      `json:"name"` //Summary
	Start       time.Time   `json:"start"`                // RFC3339 avec offset (ex: +01:00)
	End         time.Time   `json:"end"`
	Location    string      `json:"location"`
	LastUpdate  *time.Time  `json:"lastUpdate,omitempty"` // nullable dans l’iCal
}

type Agenda struct {
    Id    *uuid.UUID `json:"id"`
    UcaId string     `json:"ucaid"`
    Name  string     `json:"name"`
}

type SchedulerPayload struct {
	AgendaID uuid.UUID 	`json:"agendaId"`
	Event    Event     	`json:"event"`
}

