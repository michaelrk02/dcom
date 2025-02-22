package dcom

type Executor func(in Unmarshaler, out Marshaler)

type Stub interface {
	Execute(method string, in Unmarshaler, out Marshaler)
}
