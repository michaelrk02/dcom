package codegen

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

type DataType struct {
	ComponentPackage Package

	IdlType  string
	BaseType string
	GoType   string
	Class    Class
	Flag     Flag
}

func NewDataType(idlType string, componentPkg Package) DataType {
	baseType := BaseTypeOf(idlType)
	return DataType{
		ComponentPackage: componentPkg,
		IdlType:          idlType,
		BaseType:         baseType,
		GoType:           GoTypeOf(idlType, componentPkg),
		Class:            ClassOf(idlType),
		Flag:             FlagOf(idlType),
	}
}

func (dt DataType) ToStructure() string {
	return fmt.Sprintf("%sToStructure%s", dt.BaseType, dt.Flag.String())
}

func (dt DataType) StructureTo() string {
	return fmt.Sprintf("StructureTo%s%s", dt.BaseType, dt.Flag.String())
}

func (dt DataType) ToObject() string {
	return fmt.Sprintf("%sToObject%s", dt.BaseType, dt.Flag.String())
}

func (dt DataType) ObjectTo() string {
	return fmt.Sprintf("ObjectTo%s%s", dt.BaseType, dt.Flag.String())
}

func (dt DataType) MarshalWrite(j *jen.Statement, m string, arg jen.Code, samePackage bool) {
	if dt.Class == ClassPrimitive {
		j.Id(m).Op(".").Id(fmt.Sprintf("Write%s%s", ToPascal(dt.BaseType), dt.Flag.String())).
			Call(arg)
	} else if dt.Class == ClassStructure {
		var fn *jen.Statement
		if samePackage {
			fn = jen.Id(dt.ToStructure())
		} else {
			fn = jen.Qual(dt.ComponentPackage.Path, dt.ToStructure())
		}
		j.Id(m).Op(".").Id(fmt.Sprintf("WriteStructure%s", dt.Flag.String())).
			Call(fn.Call(arg))
	} else if dt.Class == ClassInterface {
		var fn *jen.Statement
		if samePackage {
			fn = jen.Id(dt.ToStructure())
		} else {
			fn = jen.Qual(dt.ComponentPackage.Path, dt.ToObject())
		}
		j.Id(m).Op(".").Id(fmt.Sprintf("WriteObject%s", dt.Flag.String())).
			Call(fn.Call(arg))
	}
}

func (dt DataType) UnmarshalRead(j *jen.Statement, u string, factory jen.Code, samePackage bool) {
	if dt.Class == ClassPrimitive {
		j.Id(u).Op(".").Id(fmt.Sprintf("Read%s%s", ToPascal(dt.BaseType), dt.Flag.String())).Call()
	} else if dt.Class == ClassStructure {
		var fn *jen.Statement
		if samePackage {
			fn = jen.Id(fmt.Sprintf("New%s", dt.BaseType))
		} else {
			fn = jen.Qual(dt.ComponentPackage.Path, fmt.Sprintf("New%s", dt.BaseType))
		}
		j.Id(u).Op(".").Id(fmt.Sprintf("ReadStructure%s", dt.Flag.String())).Call(fn)
	} else if dt.Class == ClassInterface {
		j.Id(u).Op(".").Id(fmt.Sprintf("ReadObject%s", dt.Flag.String())).Call(factory)
	}
}

func (dt DataType) UnmarshalConvert(j *jen.Statement, arg jen.Code, samePackage bool) jen.Code {
	if dt.Class == ClassStructure {
		if samePackage {
			return j.Id(dt.StructureTo()).Call(arg)
		} else {
			return j.Qual(dt.ComponentPackage.Path, dt.StructureTo()).Call(arg)
		}
	} else if dt.Class == ClassInterface {
		if samePackage {
			return j.Id(dt.ObjectTo()).Call(arg)
		} else {
			return j.Qual(dt.ComponentPackage.Path, dt.ObjectTo()).Call(arg)
		}
	}
	return nil
}

func (dt DataType) Zero() jen.Code {
	if dt.Class == ClassPrimitive && dt.Flag == FlagNone {
		switch dt.BaseType {
		case "bool":
			return jen.Lit(false)
		case "int":
			return jen.Lit(0)
		case "float":
			return jen.Lit(0.0)
		case "string":
			return jen.Lit("")
		}
	} else if dt.Class == ClassStructure && dt.Flag == FlagNone {
		return jen.Qual(dt.ComponentPackage.Path, dt.BaseType).Values()
	}
	return jen.Nil()
}

func (dt DataType) Render(j *jen.Statement, samePackage bool) {
	if dt.Flag == FlagOptional && dt.Class != ClassInterface {
		j.Op("*")
	} else if dt.Flag == FlagArray {
		j.Op("[]")
	}
	dt.renderGoType(j, samePackage)
}

func (dt DataType) renderGoType(j *jen.Statement, samePackage bool) {
	if dt.Class == ClassPrimitive {
		j.Id(dt.GoType)
	} else {
		goType := strings.Split(dt.GoType, "|")
		if samePackage {
			j.Id(goType[1])
		} else {
			j.Qual(goType[0], goType[1])
		}
	}
}
