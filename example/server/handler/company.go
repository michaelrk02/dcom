package handler

import (
	"strings"

	"github.com/google/uuid"
	"github.com/michaelrk02/dcom"
	"github.com/michaelrk02/dcom/example/component"
)

type Company struct {
	*dcom.ObjectHandler

	Name      string
	Metadata  component.Metadata
	Employees []component.Employee
}

func NewCompany(f dcom.Factory, instanceID *uuid.UUID) dcom.Object {
	self := &Company{
		ObjectHandler: dcom.NewObjectHandler(f, instanceID),
		Name:          "Unnamed Company",
		Employees:     make([]component.Employee, 0),
	}
	self.Object = self
	return self
}

func (self *Company) GetCLSID() uuid.UUID {
	return component.CLSIDCompany
}

func (self *Company) Dispose() {
	for _, employee := range self.Employees {
		employee.Release()
	}
}

func (self *Company) GetName() (string, error) {
	return self.Name, nil
}

func (self *Company) SetName(name string) error {
	self.Name = name
	return nil
}

func (self *Company) GetMetadata() (component.Metadata, error) {
	return self.Metadata, nil
}

func (self *Company) SetMetadata(metadata component.Metadata) error {
	self.Metadata = metadata
	return nil
}

func (self *Company) GetEmployees(keyword string, limit *int) ([]component.Employee, error) {
	resp := []component.Employee{}
	for _, emp := range self.Employees {
		if limit != nil && len(resp) >= *limit {
			break
		}

		name, err := emp.GetName()
		if err != nil {
			return nil, err
		}

		if strings.Contains(name, keyword) {
			resp = append(resp, emp)
		}
	}
	return resp, nil
}

func (self *Company) AddEmployee(employee component.Employee) error {
	employee.Acquire()
	self.Employees = append(self.Employees, employee)
	return nil
}

func (self *Company) AddEmployees(employee []component.Employee) error {
	for _, emp := range employee {
		err := self.AddEmployee(emp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Company) RemoveEmployee(employee component.Employee) error {
	newEmployees := []component.Employee{}
	for _, emp := range self.Employees {
		if emp.GetInstanceID() != employee.GetInstanceID() {
			newEmployees = append(newEmployees, emp)
		} else {
			emp.Release()
		}
	}
	self.Employees = newEmployees
	return nil
}
