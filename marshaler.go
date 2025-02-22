package dcom

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrMarshalerWrite  error = errors.New("marshaler write error")
	ErrUnmarshalerRead error = errors.New("unmarshaler read error")
)

type Marshaler interface {
	WriteBool(v bool) error
	WriteBoolArray(v []bool) error
	WriteBoolOptional(v *bool) error
	WriteError(e error) error
	WriteFloat(v float64) error
	WriteFloatArray(v []float64) error
	WriteFloatOptional(v *float64) error
	WriteInt(v int) error
	WriteIntArray(v []int) error
	WriteIntOptional(v *int) error
	WriteObject(obj Object) error
	WriteObjectArray(objs []Object) error
	WriteObjectOptional(obj Object) error
	WriteString(v string) error
	WriteStringArray(v []string) error
	WriteStringOptional(v *string) error
	WriteStructure(v Structure) error
	WriteStructureArray(v []Structure) error
	WriteStructureOptional(v Structure) error
	WriteUUID(v uuid.UUID) error
}

type Unmarshaler interface {
	ReadBool() (bool, error)
	ReadBoolArray() ([]bool, error)
	ReadBoolOptional() (*bool, error)
	ReadError() (error, error)
	ReadFloat() (float64, error)
	ReadFloatArray() ([]float64, error)
	ReadFloatOptional() (*float64, error)
	ReadInt() (int, error)
	ReadIntArray() ([]int, error)
	ReadIntOptional() (*int, error)
	ReadObject(f Factory) (Object, error)
	ReadObjectArray(f Factory) ([]Object, error)
	ReadObjectOptional(f Factory) (Object, error)
	ReadString() (string, error)
	ReadStringArray() ([]string, error)
	ReadStringOptional() (*string, error)
	ReadStructure(ref func() Structure) (Structure, error)
	ReadStructureArray(ref func() Structure) ([]Structure, error)
	ReadStructureOptional(ref func() Structure) (Structure, error)
	ReadUUID() (uuid.UUID, error)
}
