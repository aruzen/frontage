package data

import (
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

var CardTable = map[string]model.Card{}

type Card struct {
	Local string
	UUID  uuid.UUID
}
