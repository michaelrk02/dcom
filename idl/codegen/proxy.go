package codegen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

type ProxyGenerator struct {
	Interface Interface
}

func NewProxyGenerator(in Interface) *ProxyGenerator {
	return &ProxyGenerator{
		Interface: in,
	}
}

func (gen *ProxyGenerator) Generate(f *jen.File, bp *Blueprint) {
	f.Type().Id(gen.Interface.Name).Struct(
		jen.Op("*").Qual(DCOMPackage, "ObjectProxy"),
	)

	f.Func().Id(fmt.Sprintf("New%s", gen.Interface.Name)).
		Params(
			jen.Id("instanceID").Qual(UUIDPackage, "UUID"),
			jen.Id("conn").Op("*").Qual(DCOMPackage, "ProxyConnection"),
			jen.Id("f").Qual(DCOMPackage, "Factory"),
		).
		Params(jen.Qual(DCOMPackage, "Object")).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Id("&").Id(gen.Interface.Name).Values(
				jen.Dict{
					jen.Id("ObjectProxy"): jen.Qual(DCOMPackage, "NewObjectProxy").Call(jen.Id("instanceID"), jen.Id("conn"), jen.Id("f")),
				},
			))
		}).
		Line()

	f.Func().Params(jen.Id("proxy_").Op("*").Id(gen.Interface.Name)).
		Id("GetCLSID").
		Params().
		Params(jen.Qual(UUIDPackage, "UUID")).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Qual(bp.ComponentPackage.Path, fmt.Sprintf("CLSID%s", gen.Interface.Name)))
		}).
		Line()

	for _, m := range gen.Interface.Methods {
		f.Func().Params(jen.Id("proxy_").Op("*").Id(gen.Interface.Name)).
			Id(m.Name).
			ParamsFunc(func(g *jen.Group) {
				for _, p := range m.Params {
					stmt := g.Id(p.Name)
					p.DataType.Render(stmt, false)
				}
			}).
			ParamsFunc(func(g *jen.Group) {
				if m.ReturnType != "" {
					stmt := &jen.Statement{}
					m.ReturnDataType.Render(stmt, false)
					g.Add(stmt)
				}
				g.Id("error")
			}).
			BlockFunc(func(g *jen.Group) {
				if m.Access == AccessPublic {
					catchError := func(g *jen.Group, err jen.Code) {
						if m.ReturnType != "" {
							g.Return(m.ReturnDataType.Zero(), err)
						} else {
							g.Return(err)
						}
					}

					g.Var().Id("err_").Id("error")
					g.Var().Id("params_").Qual("bytes", "Buffer")
					g.Line()

					if len(m.Params) > 0 {
						g.Id("in_").Op(":=").Qual(DCOMPackage, "NewDefaultMarshaler").Call(jen.Op("&").Id("params_")).Line()
					}

					for _, p := range m.Params {
						stmt := g.Id("err_").Op("=")
						p.DataType.MarshalWrite(stmt, "in_", jen.Id(p.Name), false)
						g.If(jen.Id("err_").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
							catchError(g, jen.Id("err_"))
						}).Line()
					}

					g.List(jen.Id("resp_"), jen.Id("err_")).Op(":=").Id("proxy_").Op(".").Id("Conn").Op(".").Id("InvokeObject").Call(
						jen.Id("proxy_").Op(".").Id("GetCLSID").Call(),
						jen.Id("proxy_").Op(".").Id("GetInstanceID").Call(),
						jen.Lit(m.Name),
						jen.Op("&").Id("params_"),
					)
					g.If(jen.Id("err_").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
						catchError(g, jen.Id("err_"))
					}).Line()

					g.Id("out_").Op(":=").Qual(DCOMPackage, "NewDefaultUnmarshaler").Call(jen.Id("resp_")).Line()

					g.List(jen.Id("errRemote_"), jen.Id("err_")).Op(":=").Id("out_").Op(".").Id("ReadError").Call()
					g.If(jen.Id("err_").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
						catchError(g, jen.Id("err_"))
					})
					g.If(jen.Id("errRemote_").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
						catchError(g, jen.Id("errRemote_"))
					}).Line()

					if m.ReturnType != "" {
						factory := jen.Id("proxy_").Op(".").Id("Factory")

						stmt := g.Var().Id("vRemote_")
						m.ReturnDataType.Render(stmt, false)
						stmt.Line()

						if m.ReturnDataType.Class == ClassPrimitive {
							stmt := g.List(jen.Id("vRemote_"), jen.Id("err_")).Op("=")
							m.ReturnDataType.UnmarshalRead(stmt, "out_", factory, false)
							g.If(jen.Id("err_").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
								catchError(g, jen.Id("err_"))
							}).Line()
						} else {
							stmtRead := g.List(jen.Id("vRemoteRaw_"), jen.Id("err_")).Op(":=")
							m.ReturnDataType.UnmarshalRead(stmtRead, "out_", factory, false)
							g.If(jen.Id("err_").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
								catchError(g, jen.Id("err_"))
							})

							stmtConvert := g.Id("vRemote_").Op("=")
							m.ReturnDataType.UnmarshalConvert(stmtConvert, jen.Id("vRemoteRaw_"), false)
							stmtConvert.Line()
						}

						g.Return(jen.Id("vRemote_"), jen.Id("errRemote_"))
					} else {
						g.Return(jen.Id("errRemote_"))
					}
				} else if m.Access == AccessPrivate {
					g.Panic(jen.Lit("method not exposed"))
				}
			}).
			Line()
	}
}
