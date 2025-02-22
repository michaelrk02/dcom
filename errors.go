package dcom

import (
	"strings"
)

func Assert(err error) {
	if err != nil {
		panic(err)
	}
}

func Describe(err error) string {
	sequence := []string{}
	unwrapper := err.(interface {
		Unwrap() []error
	})
	for _, e := range unwrapper.Unwrap() {
		sequence = append(sequence, e.Error())
	}
	return strings.Join(sequence, ": ")
}
