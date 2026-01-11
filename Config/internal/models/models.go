package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Agenda struct {
	Id    *uuid.UUID `json:"id"`
	Name  string     `json:"name"`
	UcaId string     `json:"ucaid"`
}

type Alerte struct {
	Id       *uuid.UUID `json:"id"`
	Email    string     `json:"email"`
	AgendaId string     `json:"agendaid"`
}

type Event struct {
	ID          uuid.UUID   `json:"id,omitempty"`
	AgendaIDs   []uuid.UUID `json:"agendaIds"`
	UID         string      `json:"uid"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Start       time.Time   `json:"start"`
	End         time.Time   `json:"end"`
	Location    string      `json:"location"`
	LastUpdate  *time.Time  `json:"lastUpdate,omitempty"`
}

type FieldChange struct {
	Field  string `json:"field"`
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

type TimetableAlert struct {
	AgendaID uuid.UUID     `json:"agendaId"`
	UID      string        `json:"uid"`
	Changes  []FieldChange `json:"changes"`
	After    Event         `json:"after"`
	Before   *Event        `json:"before,omitempty"`
}
