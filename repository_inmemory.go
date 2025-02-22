package dcom

import (
	"errors"

	"github.com/google/uuid"
)

type inMemoryRepositoryClassMap map[uuid.UUID]inMemoryRepositoryObjectMap
type inMemoryRepositoryObjectMap map[uuid.UUID]Object

type inMemoryRepository struct {
	classMap inMemoryRepositoryClassMap
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		classMap: make(inMemoryRepositoryClassMap),
	}
}

func (r *inMemoryRepository) RegisterObject(clsid, instanceID uuid.UUID, obj Object) {
	_, ok := r.classMap[clsid]
	if !ok {
		r.classMap[clsid] = make(inMemoryRepositoryObjectMap)
	}
	r.classMap[clsid][instanceID] = obj
}

func (r *inMemoryRepository) RevokeObject(clsid, instanceID uuid.UUID) {
	_, ok := r.classMap[clsid]
	if ok {
		delete(r.classMap[clsid], instanceID)
	}
}

func (r *inMemoryRepository) ResolveObject(clsid, instanceID uuid.UUID) (Object, error) {
	_, ok := r.classMap[clsid]
	if !ok {
		return nil, errors.Join(ErrRegistryResolve, ErrRegistryResolveNotFound)
	}
	obj, ok := r.classMap[clsid][instanceID]
	if !ok {
		return nil, errors.Join(ErrRegistryResolve, ErrRegistryResolveNotFound)
	}
	return obj, nil
}
