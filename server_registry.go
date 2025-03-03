package dcom

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

var (
	ErrRegistryServerCreate    error = errors.New("server create instance error")
	ErrRegistryResolve         error = errors.New("registry resolve object error")
	ErrRegistryResolveNotFound error = errors.New("object not found")
)

type HandlerEntryPoint func(f Factory, instanceID *uuid.UUID) Object
type StubEntryPoint func(f Factory, h Object) Stub

type ServerRegistry struct {
	handlerMap map[uuid.UUID]HandlerEntryPoint
	stubMap    map[uuid.UUID]StubEntryPoint
	repo       Repository
	l          *log.Logger
}

func NewServerRegistry(repo Repository, l *log.Logger) *ServerRegistry {
	return &ServerRegistry{
		handlerMap: make(map[uuid.UUID]HandlerEntryPoint),
		stubMap:    make(map[uuid.UUID]StubEntryPoint),
		repo:       repo,
		l:          l,
	}
}

func (r *ServerRegistry) AddStub(clsid uuid.UUID, stub StubEntryPoint, handler HandlerEntryPoint) {
	r.stubMap[clsid] = stub
	r.handlerMap[clsid] = handler
}

func (r *ServerRegistry) CreateInstance(clsid uuid.UUID, instanceID *uuid.UUID) (Object, error) {
	if instanceID != nil {
		obj, err := r.ResolveObject(clsid, *instanceID)
		if err != nil && !errors.Is(err, ErrRegistryResolveNotFound) {
			return nil, errors.Join(ErrRegistryServerCreate, err)
		} else if !errors.Is(err, ErrRegistryResolveNotFound) {
			return obj, nil
		}
	}

	handler, ok := r.handlerMap[clsid]
	if !ok {
		return nil, errors.Join(ErrRegistryServerCreate, errors.New("CLSID is not registered"))
	}

	obj := handler(r, instanceID)
	r.RegisterObject(clsid, obj.GetInstanceID(), obj)

	r.l.Printf("instance created: clsid=%s id=%s", obj.GetCLSID(), obj.GetInstanceID())

	return obj, nil
}

func (r *ServerRegistry) Destroy(obj Object) {
	r.RevokeObject(obj.GetCLSID(), obj.GetInstanceID())
	r.l.Printf("instance destroyed: clsid=%s id=%s", obj.GetCLSID(), obj.GetInstanceID())
}

func (r *ServerRegistry) CreateStub(obj Object) (Stub, error) {
	stub, ok := r.stubMap[obj.GetCLSID()]
	if !ok {
		return nil, errors.Join(ErrRegistryServerCreate, errors.New("CLSID is not registered"))
	}
	return stub(r, obj), nil
}

func (r *ServerRegistry) RegisterObject(clsid, instanceID uuid.UUID, obj Object) {
	r.repo.RegisterObject(clsid, instanceID, obj)
}

func (r *ServerRegistry) RevokeObject(clsid, instanceID uuid.UUID) {
	r.repo.RevokeObject(clsid, instanceID)
}

func (r *ServerRegistry) ResolveObject(clsid, instanceID uuid.UUID) (Object, error) {
	return r.repo.ResolveObject(clsid, instanceID)
}
