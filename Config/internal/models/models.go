package models

import (
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
