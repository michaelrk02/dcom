package dcom

import (
	"bytes"

	"github.com/google/uuid"
)

type ObjectProxy struct {
	Object

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

func (self *ObjectProxy) Acquire() {
	_, err := self.Conn.InvokeObject(
		self.GetCLSID(),
		self.GetInstanceID(),
		"Acquire",
		&bytes.Buffer{},
	)
	Assert(err)
}

func (self *ObjectProxy) Release() {
	_, err := self.Conn.InvokeObject(
		self.GetCLSID(),
		self.GetInstanceID(),
		"Release",
		&bytes.Buffer{},
	)
	Assert(err)
}

func (self *ObjectProxy) Dispose() {
	panic("object disposal is not allowed in client context")
}
