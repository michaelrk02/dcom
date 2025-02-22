package dcom

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrRegistryClientCreate error = errors.New("client create instance error")
)

type ProxyEntryPoint func(instanceID uuid.UUID, conn *ProxyConnection, f Factory) Object

type ClientRegistry struct {
	conn *ProxyConnection

	proxyMap map[uuid.UUID]ProxyEntryPoint
}

func NewClientRegistry(conn *ProxyConnection) *ClientRegistry {
	return &ClientRegistry{
		conn:     conn,
		proxyMap: make(map[uuid.UUID]ProxyEntryPoint),
	}
}

func (r *ClientRegistry) AddProxy(clsid uuid.UUID, proxy ProxyEntryPoint) {
	r.proxyMap[clsid] = proxy
}

func (r *ClientRegistry) CreateInstance(clsid uuid.UUID, instanceID *uuid.UUID) (Object, error) {
	if instanceID == nil {
		return nil, errors.Join(ErrRegistryClientCreate, errors.New("client instance ID cannot be nil"))
	}

	proxy, ok := r.proxyMap[clsid]
	if !ok {
		return nil, errors.Join(ErrRegistryClientCreate, errors.New("CLSID is not registered"))
	}

	return proxy(*instanceID, r.conn, r), nil
}
