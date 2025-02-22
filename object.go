package dcom

import (
	"github.com/google/uuid"
)

type Object interface {
	GetCLSID() uuid.UUID
	GetInstanceID() uuid.UUID
}
