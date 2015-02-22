package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"
)

func TestGenerateTags(t *testing.T) {
	testCode, _ := ioutil.ReadFile("testfile.go")
	defer ioutil.WriteFile("testfile.go", testCode, 0644)
	GenerateTags("testfile.go", "json")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "testfile.go", nil, parser.ParseComments)
	if err != nil {
		t.Errorf("Error parsing generated file", err)
		return
	}

	for _, d := range f.Scope.Objects {
		if d.Kind == ast.Typ {
			ts, ok := d.Decl.(*ast.TypeSpec)
			if !ok {
				t.Errorf("Unknown type without TypeSec: %v", d)
				return
			}

			x, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range x.Fields.List {
				if len(field.Names) == 0 {
					if field.Tag != nil {
						t.Errorf("Embedded struct shouldn't be added a tag - %s", field.Tag.Value)
					}
					continue
				}
				name := field.Names[0].String()
				if name == "Field1" {
					if field.Tag == nil {
						t.Error("Tag should exist for Field1")
					} else if field.Tag.Value != "`json:\"-\"`" {
						t.Error("Shouldn't overwrite existing tags")
					}
				} else if name == "TestField2" {
					if field.Tag == nil {
						t.Error("Tag should be generated for TestFiled2")
					} else if field.Tag.Value != "`json:\"test_field2\"`" {
						t.Error("Snake case tag should be generated for TestField2")
					}
				}

			}
		}
	}
}
