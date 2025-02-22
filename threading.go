package dcom

type ThreadingModel int

const (
	ThreadingModelSingle ThreadingModel = iota
	ThreadingModelMultiple
)
