package handler

import (
	"github.com/google/uuid"
	"github.com/michaelrk02/dcom"
	"github.com/michaelrk02/dcom/example/component"
)

type Employee struct {
	*dcom.ObjectHandler

	Company component.Company

	Name    string
	Salary  float64
	Tenure  int
	Married bool
}

func NewEmployee(f dcom.Factory, instanceID *uuid.UUID) dcom.Object {
	self := &Employee{
		ObjectHandler: dcom.NewObjectHandler(f, instanceID),
		Name:          "Unnamed Employee",
	}
	self.Object = self
	return self
}

func (self *Employee) GetCLSID() uuid.UUID {
	return component.CLSIDEmployee
}

func (self *Employee) Dispose() {
}

func (self *Employee) GetCompany() (component.Company, error) {
	return self.Company, nil
}

func (self *Employee) SetCompany(company component.Company) error {
	self.Company = company
	return nil
}

func (self *Employee) GetName() (string, error) {
	return self.Name, nil
}

func (self *Employee) SetName(name string) error {
	self.Name = name
	return nil
}

func (self *Employee) GetSalary() (float64, error) {
	return self.Salary, nil
}

func (self *Employee) SetSalary(salary float64) error {
	self.Salary = salary
	return nil
}

func (self *Employee) GetTenure() (int, error) {
	return self.Tenure, nil
}

func (self *Employee) SetTenure(tenure int) error {
	self.Tenure = tenure
	return nil
}

func (self *Employee) IsMarried() (bool, error) {
	return self.Married, nil
}

func (self *Employee) SetMarried(married bool) error {
	self.Married = married
	return nil
}
