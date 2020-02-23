package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

const defaultTag = "json"
const defaultCase = "snake"
const cmdUsage = `
Usage : easytags [options] <file_name> [<tag:case>]
Examples:
- Will add json in camel case and xml in default case (snake) tags to struct fields
	easytags file.go json:camel,xml
- Will remove all tags when -r flag used when no flags provided
	easytag -r file.go
Options:

	-r removes all tags if none was provided`

type TagOpt struct {
	Tag  string
	Case string
}

func main() {
	remove := flag.Bool("r", false, "removes all tags if none was provided")
	flag.Parse()

	args := flag.Args()
	var tags []*TagOpt

	if len(args) < 1 {
		fmt.Println(cmdUsage)
		return
	} else if len(args) == 2 {
		provided := strings.Split(args[1], ",")
		for _, e := range provided {
			t := strings.SplitN(strings.TrimSpace(e), ":", 2)
			tag := &TagOpt{t[0], defaultCase}
			if len(t) == 2 {
				tag.Case = t[1]
			}
			tags = append(tags, tag)
		}
	}

	if len(tags) == 0 && *remove == false {
		tags = append(tags, &TagOpt{defaultTag, defaultCase})
	}
	for _, arg := range args {
		files, err := filepath.Glob(arg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		for _, f := range files {
			GenerateTags(f, tags, *remove)
		}
	}
}

// GenerateTags generates snake case json tags so that you won't need to write them. Can be also extended to xml or sql tags
func GenerateTags(fileName string, tags []*TagOpt, remove bool) {
	fSet := token.NewFileSet() // positions are relative to fSet
	// Parse the file given in arguments
	f, err := parser.ParseFile(fSet, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %v", err)
		return
	}

	// Range over the objects in the scope of this generated AST and check for StructType. Then range over fields
	// contained in that struct.
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.StructType:
			processTags(t, tags, remove)
			return false
		}
		return true
	})

	// Overwrite the file with modified version of ast.
	write, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error opening file %v", err)
		return
	}
	defer func() {
		err := write.Close()
		if err != nil {
			fmt.Printf("Errror writing file %v", err)
		}
	}()
	w := bufio.NewWriter(write)
	err = format.Node(w, fSet, f)
	if err != nil {
		fmt.Printf("Error formating file %s", err)
		return
	}
	defer func() {
		err := w.Flush()
		if err != nil {
			fmt.Printf("Error writing file %v", err)
		}
	}()
}

func parseTags(field *ast.Field, tags []*TagOpt) string {
	var tagValues []string
	fieldName := field.Names[0].String()
	for _, tag := range tags {
		var value string
		existingTagReg := regexp.MustCompile(fmt.Sprintf("%s:\"[^\"]+\"", tag.Tag))
		existingTag := existingTagReg.FindString(field.Tag.Value)
		if existingTag == "" {
			tName := strings.ToLower(tag.Tag)
			var name string
			switch {
			case tName == "json" || tName == "xml":
				switch tag.Case {
				case "snake":
					name = ToSnake(fieldName)
				case "camel":
					name = ToCamel(fieldName)
				case "pascal":
					name = fieldName
				default:
					fmt.Printf("Unknown case option %s", tag.Case)
				}

			case tName == "swaggertype":
				fType := strings.ToLower(fmt.Sprint(field.Type))
				switch {
				case strings.Contains(fType, "int"):
					name = "integer"
				case strings.Contains(fType, "string"):
					name = "string"
				case strings.Contains(fType, "time"):
					name = "number"
				case strings.Contains(fType, "float"):
					name = "number"
				case strings.Contains(fType, "bool"):
					name = "boolean"
				}

			}
			value = fmt.Sprintf("%s:\"%s\"", tag.Tag, name)
			tagValues = append(tagValues, value)
		}
	}
	updatedTags := strings.Fields(strings.Trim(field.Tag.Value, "`"))

	if len(tagValues) > 0 {
		updatedTags = append(updatedTags, tagValues...)
	}
	newValue := "`" + strings.Join(updatedTags, " ") + "`"

	return newValue
}

func processTags(x *ast.StructType, tags []*TagOpt, remove bool) {
	for _, field := range x.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		if !unicode.IsUpper(rune(field.Names[0].String()[0])) {
			// not exported
			continue
		}
		if remove {
			field.Tag = nil
			continue
		}
		if field.Tag == nil {
			field.Tag = &ast.BasicLit{}
			field.Tag.ValuePos = field.Type.Pos() + 1
			field.Tag.Kind = token.STRING
		}

		newTags := parseTags(field, tags)
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

// ToLowerCamel convert the given string to camelCase
func ToCamel(in string) string {
	runes := []rune(in)
	length := len(runes)

	var i int
	for i = 0; i < length; i++ {
		if unicode.IsLower(runes[i]) {
			break
		}
		runes[i] = unicode.ToLower(runes[i])
	}
	if i != 1 && i != length {
		i--
		runes[i] = unicode.ToUpper(runes[i])
	}
	return string(runes)
}
