// This file is automatically generated using DCOM IDL
// Please do not edit by hand

package component

import (
	uuid "github.com/google/uuid"
	dcom "github.com/michaelrk02/dcom"
)

var CLSIDCompany = uuid.MustParse("e1192b79-c05b-4ec5-bae4-cb6acdd9f9a0")

type Company interface {
	dcom.Object

	GetEmployees(keyword string, limit *int) ([]Employee, error)
	AddEmployee(employee Employee) error
	AddEmployees(employee []Employee) error
	RemoveEmployee(employee Employee) error
	GetName() (string, error)
	SetName(name string) error
	GetMetadata() (Metadata, error)
	SetMetadata(metadata Metadata) error
}

func CompanyToObject(v Company) dcom.Object {
	return v.(dcom.Object)
}

func CompanyToObjectOptional(v Company) dcom.Object {
	if v == nil {
		return nil
	}
	return v.(dcom.Object)
}

func CompanyToObjectArray(v []Company) []dcom.Object {
	arr := make([]dcom.Object, len(v))
	for i := range v {
		arr[i] = CompanyToObject(v[i])
	}
	return arr
}

func ObjectToCompany(v dcom.Object) Company {
	return v.(Company)
}

func ObjectToCompanyOptional(v dcom.Object) Company {
	if v == nil {
		return nil
	}
	return v.(Company)
}

func ObjectToCompanyArray(v []dcom.Object) []Company {
	arr := make([]Company, len(v))
	for i := range v {
		arr[i] = ObjectToCompany(v[i])
	}
	return arr
}
