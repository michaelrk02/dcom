// This file is automatically generated using DCOM IDL
// Please do not edit by hand

package component

import (
	dcom "github.com/michaelrk02/dcom"
	"reflect"
)

type Contact struct {
	Email     string
	Telephone string
}

func NewContact() dcom.Structure {
	return &Contact{}
}

func (s *Contact) Marshal(m dcom.Marshaler) error {
	var err error

	err = m.WriteString(s.Email)
	if err != nil {
		return err
	}

	err = m.WriteString(s.Telephone)
	if err != nil {
		return err
	}

	return nil
}

func (s *Contact) Unmarshal(u dcom.Unmarshaler) error {
	var err error

	s.Email, err = u.ReadString()
	if err != nil {
		return err
	}

	s.Telephone, err = u.ReadString()
	if err != nil {
		return err
	}

	return nil
}

func ContactToStructure(v Contact) dcom.Structure {
	return &v
}

func ContactToStructureOptional(v *Contact) dcom.Structure {
	return v
}

func ContactToStructureArray(v []Contact) []dcom.Structure {
	arr := make([]dcom.Structure, len(v))
	for i := range v {
		arr[i] = ContactToStructure(v[i])
	}
	return arr
}

func StructureToContact(v dcom.Structure) Contact {
	return *v.(*Contact)
}

func StructureToContactOptional(v dcom.Structure) *Contact {
	if reflect.ValueOf(v).IsNil() {
		return nil
	}
	return v.(*Contact)
}

func StructureToContactArray(v []dcom.Structure) []Contact {
	arr := make([]Contact, len(v))
	for i := range v {
		arr[i] = StructureToContact(v[i])
	}
	return arr
}
