package dcom

import (
	"sync/atomic"

	"github.com/google/uuid"
)

type ObjectHandler struct {
	Object

	instanceID uuid.UUID
	refs       atomic.Int64

	Factory Factory
}

func NewObjectHandler(f Factory, instanceID *uuid.UUID) *ObjectHandler {
	var instanceIDVal uuid.UUID
	if instanceID != nil {
		instanceIDVal = *instanceID
	} else {
		instanceIDVal = uuid.New()
	}

	h := &ObjectHandler{
		instanceID: instanceIDVal,
		Factory:    f,
	}
	h.refs.Store(1)

	return h
}

func (self *ObjectHandler) GetInstanceID() uuid.UUID {
	return self.instanceID
}

func (self *ObjectHandler) Acquire() {
	self.refs.Add(1)
}

func (self *ObjectHandler) Release() {
	refs := self.refs.Add(-1)
	if refs == 0 {
		self.Dispose()
		self.Factory.Destroy(self)
	}
}
