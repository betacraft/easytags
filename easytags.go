package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"unicode"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Usage : easytags {file_name} {tag_name} \n example: easytags file.go json")
		return
	}

	GenerateTags(args[0], args[1])
}

// generates snake case json tags so that you won't need to write them. Can be also exteded to xml or sql tags
func GenerateTags(fileName, tagName string) {
	fset := token.NewFileSet() // positions are relative to fset
	// Parse the file given in arguments
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return
	}

	// range over the objects in the scope of this generated AST and check for StructType. Then range over fields
	// contained in that struct.
	for _, d := range f.Scope.Objects {
		if d.Kind == ast.Typ {
			ts, ok := d.Decl.(*ast.TypeSpec)
			if !ok {
				fmt.Printf("Unknown type without TypeSec: %v", d)
				return
			}

			x, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range x.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				// if tag for field doesn't exists, create one
				if field.Tag == nil {
					name := field.Names[0].String()
					field.Tag = &ast.BasicLit{}
					field.Tag.ValuePos = field.Type.Pos() + 1
					field.Tag.Kind = token.STRING
					field.Tag.Value = fmt.Sprintf("`%s:\"%s\"`", tagName, ToSnake(name))
				} else if !strings.Contains(field.Tag.Value, fmt.Sprintf("%s:", tagName)) {
					// if tag exists, but doesn't contain target tag
					name := field.Names[0].String()
					field.Tag.Value = fmt.Sprintf("`%s:\"%s\" %s`", tagName, ToSnake(name), strings.Replace(field.Tag.Value, "`", "", 2))
				}
			}
		}
	}

	// overwrite the file with modified version of ast.
	write, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error opening file %v", err)
		return
	}
	defer write.Close()
	w := bufio.NewWriter(write)
	err = format.Node(w, fset, f)
	if err != nil {
		fmt.Printf("Error formating file", err)
		return
	}
	w.Flush()
}

// ToSnake convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
// Original source : https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}
