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
	"regexp"
)

func main() {
	args := os.Args[1:]
	var tp []string

	if len(args) < 1 {
		fmt.Println("Usage : easytags {file_name} {tag_name} \n example: easytags file.go json")
		return
	} else if len(args) == 2 {
		provided := strings.Split(args[1], ",")
		for _, e := range provided {
			tp = append(tp, strings.TrimSpace(e))
		}
	}

	GenerateTags(args[0], tp)
}

// generates snake case json tags so that you won't need to write them. Can be also exteded to xml or sql tags
func GenerateTags(fileName string, tagNames []string) {
	fset := token.NewFileSet() // positions are relative to fset
	// Parse the file given in arguments
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %v", err)
		return
	}

	// range over the objects in the scope of this generated AST and check for StructType. Then range over fields
	// contained in that struct.

	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.StructType:
			processTags(t, tagNames)
			return false
		}
		return true
	})

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
		fmt.Printf("Error formating file %s", err)
		return
	}
	w.Flush()
}

func parseTags(field *ast.Field, tags []string) string {
	var tagValues []string
	fieldName := field.Names[0].String()

	for _, tag := range tags {
		var value string
		existingTagReg := regexp.MustCompile(fmt.Sprintf("%s:\"[^\"]+\"", tag))
		existingTag := existingTagReg.FindString(field.Tag.Value)
		if existingTag != "" {
			value = existingTag
		} else {
			value = fmt.Sprintf("%s:\"%s\"", tag, ToSnake(fieldName))
		}

		tagValues = append(tagValues, value)
	}

	if len(tagValues) == 0 {
		return ""
	}

	newValue := "`" + strings.Join(tagValues, " ") + "`"

	return newValue
}

func processTags(x *ast.StructType, tagNames []string) {
	for _, field := range x.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		if field.Tag == nil {
			field.Tag = &ast.BasicLit{}
			field.Tag.ValuePos = field.Type.Pos() + 1
			field.Tag.Kind = token.STRING
		}

		newTags := parseTags(field, tagNames)
		field.Tag.Value = newTags
	}
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
