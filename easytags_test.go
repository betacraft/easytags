package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"
)

func TestGenerateTags(t *testing.T) {
	testCode, err := ioutil.ReadFile("testfile.go")
	if err != nil {
		t.Errorf("Error reading file %v", err)
	}
	defer ioutil.WriteFile("testfile.go", testCode, 0644)
	GenerateTags("testfile.go", []*TagOpt{&TagOpt{"json", "snake"}}, false)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "testfile.go", nil, parser.ParseComments)
	if err != nil {
		t.Errorf("Error parsing generated file %v", err)
		genFile, _ := ioutil.ReadFile("testfile.go")
		t.Errorf("\n%s", genFile)
		return
	}

	for _, d := range f.Scope.Objects {
		if d.Kind != ast.Typ {
			continue
		}
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
			} else if name == "ExistingTag" {
				if field.Tag == nil {
					t.Error("Tag should be generated for TestFiled2")
				} else if field.Tag.Value != "`custom:\"\" json:\"etag\"`" {
					t.Error("existing tag should not be modified, instead found ", field.Tag.Value)
				}
			}

		}
	}
}

func TestGenerateTags_Multiple(t *testing.T) {
	testCode, err := ioutil.ReadFile("testfile.go")
	if err != nil {
		t.Errorf("Error reading file %v", err)
	}
	defer ioutil.WriteFile("testfile.go", testCode, 0644)
	GenerateTags("testfile.go", []*TagOpt{&TagOpt{"json", "snake"}, &TagOpt{"xml", "snake"}}, false)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "testfile.go", nil, parser.ParseComments)
	if err != nil {
		t.Errorf("Error parsing generated file %v", err)
		return
	}

	for _, d := range f.Scope.Objects {
		if d.Kind != ast.Typ {
			continue
		}
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
				} else if field.Tag.Value != "`json:\"-\" xml:\"field1\"`" {
					t.Error("Shouldn't overwrite existing json tag, and should add xml tag")
				}
			} else if name == "TestField2" {
				if field.Tag == nil {
					t.Error("Tag should be generated for TestFiled2")
				} else if field.Tag.Value != "`json:\"test_field2\" xml:\"test_field2\"`" {
					t.Error("Snake case tag should be generated for TestField2")
				}
			} else if name == "ExistingTag" {
				if field.Tag == nil {
					t.Error("Tag should be generated for TestFiled2")
				} else if field.Tag.Value != "`custom:\"\" json:\"etag\" xml:\"existing_tag\"`" {
					t.Error("new tag should be appended to existing tag, instead found ", field.Tag.Value)
				}
			}

		}
	}
}

func TestGenerateTags_PascalCase(t *testing.T) {
	testCode, err := ioutil.ReadFile("testfile.go")
	if err != nil {
		t.Errorf("Error reading file %v", err)
	}
	defer ioutil.WriteFile("testfile.go", testCode, 0644)
	GenerateTags("testfile.go", []*TagOpt{&TagOpt{"json", "camel"}}, false)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "testfile.go", nil, parser.ParseComments)
	if err != nil {
		t.Errorf("Error parsing generated file %v", err)
		genFile, _ := ioutil.ReadFile("testfile.go")
		t.Errorf("\n%s", genFile)
		return
	}

	for _, d := range f.Scope.Objects {
		if d.Kind != ast.Typ {
			continue
		}
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
			if name == "TestField2" {
				if field.Tag == nil {
					t.Error("Tag should be generated for TestFiled2")
				} else if field.Tag.Value != "`json:\"testField2\"`" {
					t.Error("Camel tag should be generated for TestField2")
				}
			}
		}
	}
}

func TestGenerateTags_RemoveAll(t *testing.T) {
	testCode, err := ioutil.ReadFile("testfile.go")
	if err != nil {
		t.Errorf("Error reading file %v", err)
	}
	defer ioutil.WriteFile("testfile.go", testCode, 0644)
	GenerateTags("testfile.go", []*TagOpt{}, true)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "testfile.go", nil, parser.ParseComments)
	if err != nil {
		t.Errorf("Error parsing generated file %v", err)
		return
	}

	for _, d := range f.Scope.Objects {
		if d.Kind != ast.Typ {
			continue
		}
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
				if field.Tag != nil {
					t.Error("Field1 should not have any tag")
				}
			} else if name == "TestField2" {
				if field.Tag != nil {
					t.Error("TestField2 should not have any tag")
				}
			}
		}
	}
}

func TestToSnake(t *testing.T) {
	test := func(in, out string) {
		r := ToSnake(in)
		if r != out {
			t.Errorf("%s in snake_case should be %s, instead found %s", in, out, r)
		}
	}
	test("A", "a")
	test("ID", "id")
	test("UserID", "user_id")
	test("CSRFToken", "csrf_token")
}

func TestToCamel(t *testing.T) {
	test := func(in, out string) {
		r := ToCamel(in)
		if r != out {
			t.Errorf("%s in lowerCamelCase should be %s, instead found %s", in, out, r)
		}
	}
	test("A", "a")
	test("ID", "id")
	test("UserID", "userID")
	test("CSRFToken", "csrfToken")
}
