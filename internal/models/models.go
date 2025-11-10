package models

import (
	"time"
	"github.com/gofrs/uuid"
)

type Agenda struct {
	Id    *uuid.UUID `json:"id"`
	Name  string     `json:"name"`
	UcaId *uuid.UUID `json:"ucaid"`
}

type Alerte struct {
	Id       *uuid.UUID `json:"id"`
	Email    string     `json:"email"`
	AgendaId *uuid.UUID `json:"agendaid"`
}

type Event struct {
	ID          uuid.UUID   `json:"id"`
	AgendaIDs   []uuid.UUID `json:"agendaIds"`            // tableau d’UUIDs
	UID         string      `json:"uid"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Start       time.Time   `json:"start"`                // RFC3339 avec offset (ex: +01:00)
	End         time.Time   `json:"end"`
	Location    string      `json:"location"`
	LastUpdate  *time.Time  `json:"lastUpdate,omitempty"` // nullable dans l’iCal
}