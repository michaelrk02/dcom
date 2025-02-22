package dcom

import (
	"github.com/google/uuid"
)

type ObjectProxy struct {
	instanceID uuid.UUID

	Conn    *ProxyConnection
	Factory Factory
}

func NewObjectProxy(instanceID uuid.UUID, conn *ProxyConnection, f Factory) *ObjectProxy {
	return &ObjectProxy{
		instanceID: instanceID,
		Conn:       conn,
		Factory:    f,
	}
}

func (self *ObjectProxy) GetInstanceID() uuid.UUID {
	return self.instanceID
}
