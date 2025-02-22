package dcom

import "errors"

type ObjectStub struct {
	executorMap map[string]Executor
	base        Object

	Factory Factory
}

func NewObjectStub(f Factory, base Object) *ObjectStub {
	stub := &ObjectStub{
		executorMap: make(map[string]Executor),
		base:        base,
		Factory:     f,
	}

	stub.AddExecutor("Acquire", stub.ExecuteAcquire)
	stub.AddExecutor("Release", stub.ExecuteRelease)

	return stub
}

func (self *ObjectStub) AddExecutor(method string, executor Executor) {
	self.executorMap[method] = executor
}

func (self *ObjectStub) Execute(method string, in Unmarshaler, out Marshaler) {
	executor, ok := self.executorMap[method]
	if !ok {
		_ = out.WriteError(errors.New("method does not exist"))
		return
	}
	executor(in, out)
}

func (self *ObjectStub) ExecuteAcquire(in Unmarshaler, out Marshaler) {
	self.base.Acquire()
}

func (self *ObjectStub) ExecuteRelease(in Unmarshaler, out Marshaler) {
	self.base.Release()
}
