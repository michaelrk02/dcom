package dcom

import "github.com/google/uuid"

type Repository interface {
	RegisterObject(clsid, instanceID uuid.UUID, obj Object)
	RevokeObject(clsid, instanceID uuid.UUID)
	ResolveObject(clsid, instanceID uuid.UUID) (Object, error)
}
