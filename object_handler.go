package dcom

import (
	"github.com/google/uuid"
)

type ObjectHandler struct {
	instanceID uuid.UUID

	Factory Factory
}

func NewObjectHandler(f Factory, instanceID *uuid.UUID) *ObjectHandler {
	var instanceIDVal uuid.UUID
	if instanceID != nil {
		instanceIDVal = *instanceID
	} else {
		instanceIDVal = uuid.New()
	}

	return &ObjectHandler{
		instanceID: instanceIDVal,
		Factory:    f,
	}
}

func (self *ObjectHandler) GetInstanceID() uuid.UUID {
	return self.instanceID
}
