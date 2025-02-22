package codegen

import "encoding/xml"

type Access string

const (
	AccessPublic   Access = "public"
	AccessPrivate  Access = "private"
	AccessReadonly Access = "readonly"
)

type Class int

const (
	ClassPrimitive Class = iota
	ClassStructure
	ClassInterface
)

type Flag int

const (
	FlagNone Flag = iota
	FlagOptional
	FlagArray
)

func (f Flag) String() string {
	if f == FlagOptional {
		return "Optional"
	}
	if f == FlagArray {
		return "Array"
	}
	return ""
}

type Tree struct {
	XMLName xml.Name `xml:"idl"`

	Structures []Structure `xml:"structure"`
	Interfaces []Interface `xml:"interface"`
}

func (t *Tree) Merge(other *Tree) {
	t.Structures = append(t.Structures, other.Structures...)
	t.Interfaces = append(t.Interfaces, other.Interfaces...)
}

type Structure struct {
	Name       string              `xml:"name,attr"`
	Properties []StructureProperty `xml:"property"`
}

type StructureProperty struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`

	DataType DataType
}

type Interface struct {
	Name       string              `xml:"name,attr"`
	CLSID      string              `xml:"clsid,attr"`
	Properties []InterfaceProperty `xml:"property"`
	Methods    []InterfaceMethod   `xml:"method"`
}

type InterfaceProperty struct {
	Name   string `xml:"name,attr"`
	Type   string `xml:"type,attr"`
	Access Access `xml:"access,attr"`
}

type InterfaceMethod struct {
	Name       string                 `xml:"name,attr"`
	Access     Access                 `xml:"access,attr"`
	ReturnType string                 `xml:"return,attr"`
	Params     []InterfaceMethodParam `xml:"param"`

	ReturnDataType *DataType
}

type InterfaceMethodParam struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`

	DataType DataType
}
