package codegen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

type StructureGenerator struct {
	Structure Structure
}

func NewStructureGenerator(st Structure) *StructureGenerator {
	return &StructureGenerator{
		Structure: st,
	}
}

func (gen *StructureGenerator) Generate(f *jen.File, bp *Blueprint) {
	f.Type().Id(gen.Structure.Name).StructFunc(func(g *jen.Group) {
		for _, property := range gen.Structure.Properties {
			memb := g.Id(ToPascal(property.Name))
			property.DataType.Render(memb, true)
		}
	})

	f.Func().Id(fmt.Sprintf("New%s", gen.Structure.Name)).
		Params().
		Params(jen.Qual(DCOMPackage, "Structure")).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Op("&").Id(gen.Structure.Name).Values())
		}).
		Line()

	f.Func().Params(jen.Id("s").Op("*").Id(gen.Structure.Name)).Id("Marshal").
		Params(jen.Id("m").Qual(DCOMPackage, "Marshaler")).
		Params(jen.Id("error")).
		BlockFunc(func(g *jen.Group) {
			g.Var().Id("err").Id("error").Line()
			for _, prop := range gen.Structure.Properties {
				stmt := g.Id("err").Op("=")
				prop.DataType.MarshalWrite(stmt, "m", jen.Id("s").Op(".").Id(ToPascal(prop.Name)), true)
				g.If(jen.Id("err").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
					g.Return(jen.Id("err"))
				}).Line()
			}
			g.Return(jen.Nil())
		}).
		Line()

	f.Func().Params(jen.Id("s").Op("*").Id(gen.Structure.Name)).Id("Unmarshal").
		Params(jen.Id("u").Qual(DCOMPackage, "Unmarshaler")).
		Params(jen.Id("error")).
		BlockFunc(func(g *jen.Group) {
			g.Var().Id("err").Id("error").Line()
			for _, prop := range gen.Structure.Properties {
				if prop.DataType.Class == ClassPrimitive {
					stmt := g.List(jen.Id("s").Op(".").Id(ToPascal(prop.Name)), jen.Id("err")).Op("=")
					prop.DataType.UnmarshalRead(stmt, "u", jen.Nil(), true)
					g.If(jen.Id("err").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
						g.Return(jen.Id("err"))
					}).Line()
				} else if prop.DataType.Class == ClassStructure {
					stmtRead := g.List(jen.Id(prop.Name), jen.Id("err")).Op(":=")
					prop.DataType.UnmarshalRead(stmtRead, "u", jen.Nil(), true)

					g.If(jen.Id("err").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
						g.Return(jen.Id("err"))
					})

					stmtConvert := g.Id("s").Op(".").Id(ToPascal(prop.Name)).Op("=")
					prop.DataType.UnmarshalConvert(stmtConvert, jen.Id(prop.Name), true)
					stmtConvert.Line()
				}
			}
			g.Return(jen.Nil())
		}).
		Line()

	f.Func().Id(fmt.Sprintf("%sToStructure", gen.Structure.Name)).
		Params(jen.Id("v").Id(gen.Structure.Name)).
		Params(jen.Qual(DCOMPackage, "Structure")).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Op("&").Id("v"))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("%sToStructureOptional", gen.Structure.Name)).
		Params(jen.Id("v").Op("*").Id(gen.Structure.Name)).
		Params(jen.Qual(DCOMPackage, "Structure")).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Id("v"))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("%sToStructureArray", gen.Structure.Name)).
		Params(jen.Id("v").Op("[]").Id(gen.Structure.Name)).
		Params(jen.Op("[]").Qual(DCOMPackage, "Structure")).
		BlockFunc(func(g *jen.Group) {
			g.Id("arr").Op(":=").Make(jen.Op("[]").Qual(DCOMPackage, "Structure"), jen.Len(jen.Id("v")))
			g.For(jen.Id("i").Op(":=").Range().Id("v")).BlockFunc(func(g *jen.Group) {
				g.Id("arr").Index(jen.Id("i")).Op("=").Id(fmt.Sprintf("%sToStructure", gen.Structure.Name)).Call(jen.Id("v").Index(jen.Id("i")))
			})
			g.Return(jen.Id("arr"))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("StructureTo%s", gen.Structure.Name)).
		Params(jen.Id("v").Qual(DCOMPackage, "Structure")).
		Params(jen.Id(gen.Structure.Name)).
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Op("*").Id("v").Op(".").Parens(jen.Op("*").Id(gen.Structure.Name)))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("StructureTo%sOptional", gen.Structure.Name)).
		Params(jen.Id("v").Qual(DCOMPackage, "Structure")).
		Params(jen.Op("*").Id(gen.Structure.Name)).
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Qual("reflect", "ValueOf").Call(jen.Id("v")).Op(".").Id("IsNil").Call()).BlockFunc(func(g *jen.Group) {
				g.Return(jen.Nil())
			})
			g.Return(jen.Id("v").Op(".").Parens(jen.Op("*").Id(gen.Structure.Name)))
		}).
		Line()

	f.Func().Id(fmt.Sprintf("StructureTo%sArray", gen.Structure.Name)).
		Params(jen.Id("v").Op("[]").Qual(DCOMPackage, "Structure")).
		Params(jen.Op("[]").Id(gen.Structure.Name)).
		BlockFunc(func(g *jen.Group) {
			g.Id("arr").Op(":=").Make(jen.Op("[]").Id(gen.Structure.Name), jen.Len(jen.Id("v")))
			g.For(jen.Id("i").Op(":=").Range().Id("v")).BlockFunc(func(g *jen.Group) {
				g.Id("arr").Index(jen.Id("i")).Op("=").Id(fmt.Sprintf("StructureTo%s", gen.Structure.Name)).Call(jen.Id("v").Index(jen.Id("i")))
			})
			g.Return(jen.Id("arr"))
		}).
		Line()
}
