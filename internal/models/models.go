package models

import (
	"github.com/gofrs/uuid"
)

type User struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
}
