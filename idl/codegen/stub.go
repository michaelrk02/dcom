package codegen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

type StubGenerator struct {
	Interface Interface
}

func NewStubGenerator(in Interface) *StubGenerator {
	return &StubGenerator{
		Interface: in,
	}
}

func (gen *StubGenerator) Generate(f *jen.File, bp *Blueprint) {
	f.Type().Id(gen.Interface.Name).Struct(
		jen.Op("*").Qual(DCOMPackage, "ObjectStub"),
		jen.Id("obj").Qual(bp.ComponentPackage.Path, gen.Interface.Name),
	)

	f.Func().Id(fmt.Sprintf("New%s", gen.Interface.Name)).
		Params(
			jen.Id("f").Qual(DCOMPackage, "Factory"),
			jen.Id("obj").Qual(DCOMPackage, "Object"),
		).
		Params(jen.Qual(DCOMPackage, "Stub")).
		BlockFunc(func(g *jen.Group) {
			g.Id("stub").Op(":=").Op("&").Id(gen.Interface.Name).Values(jen.Dict{
				jen.Id("ObjectStub"): jen.Qual(DCOMPackage, "NewObjectStub").Call(jen.Id("f")),
				jen.Id("obj"):        jen.Id("obj").Op(".").Parens(jen.Qual(bp.ComponentPackage.Path, gen.Interface.Name)),
			}).Line()

			for _, m := range gen.Interface.Methods {
				if m.Access == AccessPublic {
					g.Id("stub").Op(".").Id("AddExecutor").Call(
						jen.Lit(m.Name),
						jen.Id("stub").Op(".").Id(fmt.Sprintf("Execute%s", m.Name)),
					)
				}
			}
			g.Line()

			g.Return(jen.Id("stub"))
		}).Line()

	for _, m := range gen.Interface.Methods {
		if m.Access == AccessPublic {
			f.Func().Params(jen.Id("stub_").Op("*").Id(gen.Interface.Name)).
				Id(fmt.Sprintf("Execute%s", m.Name)).
				Params(
					jen.Id("in_").Qual(DCOMPackage, "Unmarshaler"),
					jen.Id("out_").Qual(DCOMPackage, "Marshaler"),
				).
				Params().
				BlockFunc(func(g *jen.Group) {
					factory := jen.Id("stub_").Op(".").Id("Factory")

					for _, p := range m.Params {
						if p.DataType.Class == ClassPrimitive {
							stmt := g.List(jen.Id(p.Name), jen.Id("err_")).Op(":=")
							p.DataType.UnmarshalRead(stmt, "in_", factory, false)
							g.Qual(DCOMPackage, "Assert").Call(jen.Id("err_")).Line()
						} else {
							stmtRead := g.List(jen.Id(fmt.Sprintf("%sTemp", p.Name)), jen.Id("err_")).Op(":=")
							p.DataType.UnmarshalRead(stmtRead, "in_", factory, false)
							g.Qual(DCOMPackage, "Assert").Call(jen.Id("err_"))

							stmtConvert := g.Id(p.Name).Op(":=")
							p.DataType.UnmarshalConvert(stmtConvert, jen.Id(fmt.Sprintf("%sTemp", p.Name)), false)
							stmtConvert.Line()
						}
					}

					var stmtResp *jen.Statement
					if m.ReturnType != "" {
						stmtResp = g.List(jen.Id("resp_"), jen.Id("err_")).Op(":=")
					} else {
						stmtResp = g.Id("err_").Op("=")
					}

					stmtResp.Id("stub_").Op(".").Id("obj").Op(".").Id(m.Name).ParamsFunc(func(g *jen.Group) {
						for _, p := range m.Params {
							g.Id(p.Name)
						}
					}).Line()

					g.Qual(DCOMPackage, "Assert").Call(jen.Id("out_").Op(".").Id("WriteError").Call(jen.Id("err_")))
					if m.ReturnType != "" {
						stmt := &jen.Statement{}
						m.ReturnDataType.MarshalWrite(stmt, "out_", jen.Id("resp_"), false)
						g.Qual(DCOMPackage, "Assert").Call(stmt)
					}
				}).
				Line()
		}
	}
}
