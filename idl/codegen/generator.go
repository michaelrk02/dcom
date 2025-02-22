package codegen

import "github.com/dave/jennifer/jen"

type Generator interface {
	Generate(f *jen.File, bp *Blueprint)
}
