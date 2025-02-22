package main

import (
	"fmt"
	"strings"

	"github.com/michaelrk02/dcom/idl/codegen"
)

func main() {
	blueprintFile := "blueprint.json"

	var bp codegen.Blueprint
	err := bp.Load(blueprintFile)
	if err != nil {
		panic(err)
	}

	tree, err := bp.BuildTree()
	if err != nil {
		panic(err)
	}

	for _, st := range tree.Structures {
		gen := codegen.NewStructureGenerator(st)
		dir := fmt.Sprintf("%s/component", bp.GeneratedDir)
		name := strings.ToLower(st.Name) + bp.ComponentPackage.Suffix

		err := bp.Generate(gen, dir, name, bp.ComponentPackage.Alias)
		if err != nil {
			panic(err)
		}
	}

	for _, in := range tree.Interfaces {
		gen := codegen.NewInterfaceGenerator(in)
		dir := fmt.Sprintf("%s/component", bp.GeneratedDir)
		name := strings.ToLower(in.Name) + bp.ComponentPackage.Suffix

		err := bp.Generate(gen, dir, name, bp.ComponentPackage.Alias)
		if err != nil {
			panic(err)
		}
	}

	for _, in := range tree.Interfaces {
		gen := codegen.NewProxyGenerator(in)
		dir := fmt.Sprintf("%s/proxy", bp.GeneratedDir)
		name := strings.ToLower(in.Name) + bp.ProxyPackage.Suffix

		err := bp.Generate(gen, dir, name, bp.ProxyPackage.Alias)
		if err != nil {
			panic(err)
		}
	}

	for _, in := range tree.Interfaces {
		gen := codegen.NewStubGenerator(in)
		dir := fmt.Sprintf("%s/stub", bp.GeneratedDir)
		name := strings.ToLower(in.Name) + bp.StubPackage.Suffix

		err := bp.Generate(gen, dir, name, bp.StubPackage.Alias)
		if err != nil {
			panic(err)
		}
	}
}
