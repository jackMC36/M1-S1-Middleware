package models

import (
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
