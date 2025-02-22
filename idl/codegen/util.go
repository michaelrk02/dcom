package codegen

import (
	"fmt"
	"strings"
)

func ToPascal(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func FlagOf(t string) Flag {
	if len(t) > 2 && t[len(t)-2:] == "[]" {
		return FlagArray
	}
	if len(t) > 1 && t[len(t)-1:] == "?" {
		return FlagOptional
	}
	return FlagNone
}

func BaseTypeOf(t string) string {
	t = strings.ReplaceAll(t, "[]", "")
	t = strings.ReplaceAll(t, "?", "")
	t = strings.ReplaceAll(t, "@", "")
	t = strings.ReplaceAll(t, "#", "")
	return t
}

func GoTypeOf(t string, componentPkg Package) string {
	bt := BaseTypeOf(t)

	if ClassOf(t) == ClassPrimitive {
		if bt == "bool" {
			return "bool"
		}
		if bt == "int" {
			return "int"
		}
		if bt == "float" {
			return "float64"
		}
		if bt == "string" {
			return "string"
		}
	}

	return fmt.Sprintf("%s|%s", componentPkg.Path, bt)
}

func ClassOf(t string) Class {
	if t[0] == '#' {
		return ClassStructure
	}
	if t[0] == '@' {
		return ClassInterface
	}
	return ClassPrimitive
}
