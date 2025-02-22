package dcom

import "errors"

type ObjectStub struct {
	executorMap map[string]Executor

	Factory Factory
}

func NewObjectStub(f Factory) *ObjectStub {
	return &ObjectStub{
		executorMap: make(map[string]Executor),
		Factory:     f,
	}
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
