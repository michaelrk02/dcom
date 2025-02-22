package dcom

import "github.com/google/uuid"

type Factory interface {
	CreateInstance(clsid uuid.UUID, instanceID *uuid.UUID) (Object, error)
	Destroy(obj Object)
}
