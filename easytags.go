package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

// generates snake case json tags so that you won't need to write them. Can be also exteded to xml or sql tags
func main() {
	fset := token.NewFileSet() // positions are relative to fset
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Usage : easytags {file_name} {tag_name} {debug (true/false)} \n example: easytags file.go json true")
		return
	}
	debug := false
	if len(args) == 3 {
		if args[2] == "true" {
			debug = true
		}
	}
	tagName := args[1]
	// Parse the file given in arguments
	f, err := parser.ParseFile(fset, args[0], nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return
	}
	// read entire source file as a slice of lines
	lines, err := readLines(args[0])
	if err != nil {
		fmt.Printf("Error reading file %v \n ", err)
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
				line := fset.File(field.Pos()).Line(field.Pos())
				line = line - 1
				// if tag for field doesn't exists, create one
				if field.Tag == nil {
					name := field.Names[0].String()
					if debug {
						fmt.Printf("Replacing line %s \n ", lines[line])
					}
					lines[line] = fmt.Sprintf("%s %v `%s:\"%s\"`", name, field.Type, tagName, ToSnake(name))
					if debug {
						fmt.Printf("By line : %s \n", lines[line])
					}
				} else if !strings.Contains(field.Tag.Value, fmt.Sprintf("%s:", tagName)) {
					// if tag exists, but doesn't contain target tag
					name := field.Names[0].String()
					if debug {
						fmt.Printf("Replacing line %s \n ", lines[line])
					}
					lines[line] = fmt.Sprintf("%s %v `%s:\"%s\" %s`", name, field.Type, tagName, ToSnake(name), strings.Replace(field.Tag.Value, "`", "", 2))
					if debug {
						fmt.Printf("By line : %s \n", lines[line])
					}

				}

			}
		}
	}
	// overwrite the file with modified version of lines.
	writeLines(lines, args[0])
	cmd := exec.Command("go", "fmt", args[0])
	cmd.Run()
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
// original source : http://stackoverflow.com/questions/5884154/golang-read-text-file-into-string-array-and-write
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
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
