// This file is automatically generated using DCOM IDL
// Please do not edit by hand

package stub

import (
	dcom "github.com/michaelrk02/dcom"
	component "github.com/michaelrk02/dcom/example/component"
)

type Employee struct {
	*dcom.ObjectStub
	obj component.Employee
}

func NewEmployee(f dcom.Factory, obj dcom.Object) dcom.Stub {
	stub := &Employee{
		ObjectStub: dcom.NewObjectStub(f),
		obj:        obj.(component.Employee),
	}

	stub.AddExecutor("GetName", stub.ExecuteGetName)
	stub.AddExecutor("GetSalary", stub.ExecuteGetSalary)
	stub.AddExecutor("GetTenure", stub.ExecuteGetTenure)
	stub.AddExecutor("IsMarried", stub.ExecuteIsMarried)

	return stub
}

func (stub_ *Employee) ExecuteGetName(in_ dcom.Unmarshaler, out_ dcom.Marshaler) {
	resp_, err_ := stub_.obj.GetName()

	dcom.Assert(out_.WriteError(err_))
	dcom.Assert(out_.WriteString(resp_))
}

func (stub_ *Employee) ExecuteGetSalary(in_ dcom.Unmarshaler, out_ dcom.Marshaler) {
	resp_, err_ := stub_.obj.GetSalary()

	dcom.Assert(out_.WriteError(err_))
	dcom.Assert(out_.WriteFloat(resp_))
}

func (stub_ *Employee) ExecuteGetTenure(in_ dcom.Unmarshaler, out_ dcom.Marshaler) {
	resp_, err_ := stub_.obj.GetTenure()

	dcom.Assert(out_.WriteError(err_))
	dcom.Assert(out_.WriteInt(resp_))
}

func (stub_ *Employee) ExecuteIsMarried(in_ dcom.Unmarshaler, out_ dcom.Marshaler) {
	resp_, err_ := stub_.obj.IsMarried()

	dcom.Assert(out_.WriteError(err_))
	dcom.Assert(out_.WriteBool(resp_))
}
