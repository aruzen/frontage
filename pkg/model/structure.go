package model

import "github.com/google/uuid"

type Structure interface {
	ID() uuid.UUID
}
