package codegen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

type InterfaceGenerator struct {
	Interface Interface
}

func NewInterfaceGenerator(in Interface) *InterfaceGenerator {
	return &InterfaceGenerator{
		Interface: in,
	}
}

func (gen *InterfaceGenerator) Generate(f *jen.File, bp *Blueprint) {
	f.Var().Id(fmt.Sprintf("CLSID%s", gen.Interface.Name)).Op("=").Qual(UUIDPackage, "MustParse").Call(jen.Lit(gen.Interface.CLSID))

	f.Type().Id(gen.Interface.Name).InterfaceFunc(func(g *jen.Group) {
		g.Qual(DCOMPackage, "Object")
		g.Line()

		for _, method := range gen.Interface.Methods {
			m := g.Id(method.Name).ParamsFunc(func(g *jen.Group) {
				for _, param := range method.Params {
					p := g.Id(param.Name)
					param.DataType.Render(p, true)
				}
			})
			m.ParamsFunc(func(g *jen.Group) {
				if method.ReturnDataType != nil {
					p := &jen.Statement{}
					method.ReturnDataType.Render(p, true)
					g.Add(p)
				}

				g.Id("error")
			})
		}
	})

	f.Func().Id(fmt.Sprintf("%sToObject", gen.Interface.Name)).
		Params(jen.Id("v").Id(gen.Interface.Name)).
		Params(jen.Qual(DCOMPackage, "Object")).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Id("v").Op(".").Params(jen.Qual(DCOMPackage, "Object")))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("%sToObjectOptional", gen.Interface.Name)).
		Params(jen.Id("v").Id(gen.Interface.Name)).
		Params(jen.Qual(DCOMPackage, "Object")).
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id("v").Op("==").Nil()).BlockFunc(func(g *jen.Group) {
				g.Return(jen.Nil())
			})
			g.Return(jen.Id("v").Op(".").Params(jen.Qual(DCOMPackage, "Object")))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("%sToObjectArray", gen.Interface.Name)).
		Params(jen.Id("v").Op("[]").Id(gen.Interface.Name)).
		Params(jen.Op("[]").Qual(DCOMPackage, "Object")).
		BlockFunc(func(g *jen.Group) {
			g.Id("arr").Op(":=").Make(jen.Op("[]").Qual(DCOMPackage, "Object"), jen.Len(jen.Id("v")))
			g.For(jen.Id("i").Op(":=").Range().Id("v")).BlockFunc(func(g *jen.Group) {
				g.Id("arr").Index(jen.Id("i")).Op("=").Id(fmt.Sprintf("%sToObject", gen.Interface.Name)).Call(jen.Id("v").Index(jen.Id("i")))
			})
			g.Return(jen.Id("arr"))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("ObjectTo%s", gen.Interface.Name)).
		Params(jen.Id("v").Qual(DCOMPackage, "Object")).
		Params(jen.Id(gen.Interface.Name)).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Id("v").Op(".").Params(jen.Id(gen.Interface.Name)))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("ObjectTo%sOptional", gen.Interface.Name)).
		Params(jen.Id("v").Qual(DCOMPackage, "Object")).
		Params(jen.Id(gen.Interface.Name)).
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id("v").Op("==").Nil()).BlockFunc(func(g *jen.Group) {
				g.Return(jen.Nil())
			})
			g.Return(jen.Id("v").Op(".").Params(jen.Id(gen.Interface.Name)))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("ObjectTo%sArray", gen.Interface.Name)).
		Params(jen.Id("v").Op("[]").Qual(DCOMPackage, "Object")).
		Params(jen.Op("[]").Id(gen.Interface.Name)).
		BlockFunc(func(g *jen.Group) {
			g.Id("arr").Op(":=").Make(jen.Op("[]").Id(gen.Interface.Name), jen.Len(jen.Id("v")))
			g.For(jen.Id("i").Op(":=").Range().Id("v")).BlockFunc(func(g *jen.Group) {
				g.Id("arr").Index(jen.Id("i")).Op("=").Id(fmt.Sprintf("ObjectTo%s", gen.Interface.Name)).Call(jen.Id("v").Index(jen.Id("i")))
			})
			g.Return(jen.Id("arr"))
		}).
		Line()
}
