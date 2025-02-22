package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/michaelrk02/dcom"
	"github.com/michaelrk02/dcom/example/client/proxy"
	"github.com/michaelrk02/dcom/example/component"
)

var DefaultCompanyID = uuid.MustParse("c86ef9cc-0a8d-4ab4-8068-0c37e66421f6")

func main() {
	var obj dcom.Object
	var err error

	conn := dcom.NewProxyConnection("localhost:5560")

	reg := dcom.NewClientRegistry(conn)
	reg.AddProxy(component.CLSIDCompany, proxy.NewCompany)
	reg.AddProxy(component.CLSIDEmployee, proxy.NewEmployee)

	obj, err = reg.CreateInstance(component.CLSIDCompany, &DefaultCompanyID)
	if err != nil {
		panic(err)
	}
	company := component.ObjectToCompany(obj)

	company.Acquire()
	defer company.Release()

	running := true
	for running {
		var choice int

		fmt.Println("Actions:")
		fmt.Println(" [1] View company metadata")
		fmt.Println(" [2] View list of employees")
		fmt.Println(" [0] Exit")
		fmt.Printf("Choice: ")
		fmt.Scanf("%d", &choice)

		if choice == 1 {
			metadata, err := company.GetMetadata()
			if err != nil {
				log.Println("error:", err)
				continue
			}

			fmt.Println("========================================")
			fmt.Printf("Name: %s\n", metadata.Name)
			fmt.Printf("Address: %s\n", metadata.Address)
			fmt.Printf("Website: %s\n", metadata.Website)
			fmt.Printf("Email: %s\n", metadata.Contact.Email)
			fmt.Printf("Telephone: %s\n", metadata.Contact.Telephone)
			fmt.Println("========================================")
		} else if choice == 2 {
			var keyword string
			var limitN int
			var limit *int

			fmt.Printf("Keyword: ")
			fmt.Scanf("%s", &keyword)

			fmt.Printf("Limit (0 for unlimited): ")
			fmt.Scanf("%d", &limitN)

			if limitN > 0 {
				limit = &limitN
			}

			employees, err := company.GetEmployees(keyword, limit)
			if err != nil {
				log.Println("error:", err)
				continue
			}

			for _, employee := range employees {
				employee.Acquire()

				name, err := employee.GetName()
				if err != nil {
					log.Println("error:", err)
					employee.Release()
					continue
				}

				salary, err := employee.GetSalary()
				if err != nil {
					log.Println("error:", err)
					employee.Release()
					continue
				}

				tenure, err := employee.GetTenure()
				if err != nil {
					log.Println("error:", err)
					employee.Release()
					continue
				}

				married, err := employee.GetMarried()
				if err != nil {
					log.Println("error:", err)
					employee.Release()
					continue
				}

				fmt.Println("========================================")
				fmt.Printf("Name: %s\n", name)
				fmt.Printf("Salary: %.2f\n", salary)
				fmt.Printf("Tenure: %d\n", tenure)
				fmt.Printf("Married: %t\n", married)

				employee.Release()
			}
			fmt.Println("========================================")
		} else {
			running = false
		}
	}
}
