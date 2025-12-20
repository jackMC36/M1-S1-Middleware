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

//DTStart format : "2006-01-02T15:04:05-07:00"
//DTEnd format : "2006-01-02T15:04:05-07:00"

//Summary: Option CM Méthodes approchées
//DESCRIPTION:\n\nMASTER 1 INFO\nNGUYEN MINH HIEU\n\n(Updated :20/11/2024 1
//UID:ADE60323032342d323032352d5543412d33343338392d302d34
//LAST-MODIFIED:20241120T170100Z


type SchedulerPayload struct {
	AgendaID uuid.UUID 	`json:"agendaId"`
	Event    Event     	`json:"event"`
}

type EventChange struct {
    Before *Event `json:"before,omitempty"`
	After  Event  `json:"after"`
}
