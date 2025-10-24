package models

import (
	"github.com/gofrs/uuid"
)

type Agenda struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
}
