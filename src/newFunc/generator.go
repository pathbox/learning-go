package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

var templ = `
package {{.PkgName}}
import (
{{ range .Imports -}}
    {{ if .Name -}}
    {{.Name}} {{.Path}}
    {{ else -}}
    {{.Path}}
    {{ end -}}
{{ end -}}
)
{{ range .Constructors }}
// {{.Name}} Create a new {{.Struct}}
func {{.Name}}({{.Params}}) {{if not .ValueFlag}}*{{end -}} {{.Struct}} {
	{{ if .InitFlag -}}
    s := {{if not .ValueFlag}}&{{end -}} {{.Struct}} {
        {{.Fields}}
    }
	s.init()
	return s
	{{ else -}}
    return {{if not .ValueFlag}}&{{end -}} {{.Struct}} {
        {{.Fields}}
    }
	{{ end -}}
}
{{ end }}
`

// GenerateCode generate constructors code
func GenerateCode(pkgName string, importInfos []ImportInfo, structInfos []StructInfo) (string, error) {
	// remove duplicate imports
	importInfos = filterUniqueImports(importInfos)

	// generate code with template
	t, err := template.New("").Parse(templ)
	if err != nil {
		return "", err
	}
	data := o{
		"PkgName":     pkgName,
		"Imports":     importInfos,
		"Constructor": []o{},
	}
	constructors := []o{}
	for _, structInfo := range structInfos {
		params := []string{}
		fields := []string{}
		for _, field := range structInfo.Fields {
			if field.Skipped {
				continue
			}
			params = append(params, fmt.Sprintf("%v %v", toLowerCamel(field.Name), field.Type))
			fields = append(fields, fmt.Sprintf("%v: %v,", field.Name, toLowerCamel(field.Name)))
		}
		constructors = append(constructors, o{
			"Name":      "New" + strcase.ToCamel(structInfo.StructName),
			"Struct":    structInfo.StructName,
			"InitFlag":  structInfo.InitFlag,
			"ValueFlag": structInfo.ValueFlag,
			"Params":    strings.Join(params, ", "),
			"Fields":    strings.Join(fields, "\n"),
		})
	}
	data["Constructors"] = constructors
	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	if err != nil {
		return "", err
	}

	// format code
	buf, err := formatCode(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func formatCode(source []byte) ([]byte, error) {
	output, err := imports.Process("", source, &imports.Options{
		AllErrors: true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
		Fragment:  true,
	})
	if err != nil {
		return output, err
	}
	if bytes.Equal(source, output) {
		return output, nil
	}
	return formatCode(output)
}

// filterUniqueImports remove duplicate imports, return unqiue imports
func filterUniqueImports(imports []ImportInfo) []ImportInfo {
	hash := map[string]ImportInfo{}
	for _, importInfo := range imports {
		key := fmt.Sprintf("%v|%v", importInfo.Name, importInfo.Path)
		hash[key] = importInfo
	}
	ret := []ImportInfo{}
	for _, importInfo := range hash {
		ret = append(ret, importInfo)
	}
	return ret
}

type o map[string]interface{}
