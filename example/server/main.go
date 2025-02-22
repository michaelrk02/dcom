package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/michaelrk02/dcom"
	"github.com/michaelrk02/dcom/example/component"
	"github.com/michaelrk02/dcom/example/server/handler"
	"github.com/michaelrk02/dcom/example/server/stub"
)

var DefaultCompanyID = uuid.MustParse("c86ef9cc-0a8d-4ab4-8068-0c37e66421f6")

func main() {
	logger := log.Default()

	repo := dcom.NewInMemoryRepository()

	reg := dcom.NewServerRegistry(repo)

	reg.AddHandler(component.CLSIDCompany, handler.NewCompany)
	reg.AddHandler(component.CLSIDEmployee, handler.NewEmployee)

	reg.AddStub(component.CLSIDCompany, stub.NewCompany)
	reg.AddStub(component.CLSIDEmployee, stub.NewEmployee)

	initRegistry(reg)

	conn := dcom.NewStubConnection(
		logger,
		":5560",
		dcom.ThreadingModelMultiple,
		reg,
	)
	go func() {
		err := conn.Listen()
		if err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	logger.Println("gracefully shutting down ...")
	conn.Close()
}

func initRegistry(reg *dcom.ServerRegistry) {
	var obj dcom.Object
	var err error

	obj, err = reg.CreateInstance(component.CLSIDCompany, &DefaultCompanyID)
	if err != nil {
		panic(err)
	}
	company := component.ObjectToCompany(obj)

	company.SetName("PT Mencari Cinta Sejati")
	company.SetMetadata(component.Metadata{
		Name:    "PT Mencari Cinta Sejati",
		Address: "Palur Wetan, RT 02/RW 04, Palur, Mojolaban, Sukoharjo",
		Website: "michaelrk02.my.id",
		Contact: component.Contact{
			Email:     "michaelkrisnadhi@gmail.com",
			Telephone: "(+62) 895-3438-45423",
		},
	})

	names := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	for i := range names {
		tenure := i + 1

		emp, err := createEmployee(
			reg,
			names[i],
			300.0+50.0*float64(tenure),
			tenure,
			i%2 != 0,
		)
		if err != nil {
			panic(err)
		}

		err = company.AddEmployee(emp)
		if err != nil {
			panic(err)
		}
	}
}

func createEmployee(
	reg *dcom.ServerRegistry,
	name string,
	salary float64,
	tenure int,
	married bool,
) (component.Employee, error) {
	var obj dcom.Object
	var err error

	obj, err = reg.CreateInstance(component.CLSIDEmployee, nil)
	if err != nil {
		return nil, err
	}
	employee := component.ObjectToEmployee(obj)

	employee.SetName(name)
	employee.SetSalary(salary)
	employee.SetTenure(tenure)
	employee.SetMarried(married)

	return employee, nil
}
