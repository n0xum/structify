package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Field struct {
	Name       string
	Type        string
	DatabaseTag string
}

type Struct struct {
	Name       string
	Fields      []Field
	PackageName string
	TableName   string
}

type Parser struct {
	fset    *token.FileSet
	structs  map[string][]*Struct
	pkgName string
}

type visitor struct {
	p *Parser
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		if genDecl.Tok == token.PACKAGE {
			importSpec, ok := genDecl.Specs[0].(*ast.ImportSpec)
			if ok && importSpec.Name != nil {
				v.p.pkgName = importSpec.Name.Name
			}
		}
		if genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				structName := typeSpec.Name.Name
				if !token.IsExported(structName) {
					continue
				}
				s := &Struct{
					Name:       structName,
					PackageName: v.p.pkgName,
					TableName:   toSnakeCase(structName),
				}
				v.p.extractFields(structType, s)
				if len(s.Fields) > 0 {
					if v.p.pkgName == "" {
						v.p.pkgName = "main"
					}
					v.p.structs[v.p.pkgName] = append(v.p.structs[v.p.pkgName], s)
				}
			}
		}
	}
	return v
}

func New() *Parser {
	return &Parser{
		fset:   token.NewFileSet(),
		structs: make(map[string][]*Struct),
	}
}

func (p *Parser) ParseFiles(paths []string) error {
	for _, path := range paths {
		if err := p.parseFile(path); err != nil {
			return fmt.Errorf("parse file %s: %w", path, err)
		}
	}
	return nil
}

func (p *Parser) parseFile(path string) error {
	f, err := parser.ParseFile(p.fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Walk(&visitor{p: p}, f)

	return nil
}

func (p *Parser) extractFields(structType *ast.StructType, s *Struct) {
	for _, field := range structType.Fields.List {
		if field.Names == nil {
			continue
		}

		fieldName := field.Names[0].Name
		if !token.IsExported(fieldName) && len(field.Names) == 1 {
			continue
		}

		f := Field{
			Name: fieldName,
		}

		if field.Tag != nil {
			tag := strings.Trim(field.Tag.Value, "`")
			f.DatabaseTag = parseDBTag(tag)
		}

		if field.Type != nil {
			f.Type = exprToString(field.Type)
		}

		if f.Type == "" {
			f.Type = "interface{}"
		}

		s.Fields = append(s.Fields, f)
	}
}

func (p *Parser) GetStructs() map[string][]*Struct {
	return p.structs
}

func exprToString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return exprToString(v.X)
	case *ast.StarExpr:
		return "*" + exprToString(v.X)
	case *ast.ArrayType:
		return "[]" + exprToString(v.Elt)
	case *ast.MapType:
		return "map[" + exprToString(v.Key) + "]" + exprToString(v.Value)
	case *ast.ChanType:
		return "chan"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "..."
	default:
		return "any"
	}
}

func parseDBTag(tag string) string {
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "db:") {
			value := strings.TrimPrefix(part, "db:")
			value = strings.Trim(value, `"`)
			return value
		}
	}
	return ""
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
