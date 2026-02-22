package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/n0xum/structify/internal/util"
)

type Field struct {
	Name        string
	Type        string
	DatabaseTag string
}

type Struct struct {
	Name        string
	Fields      []Field
	PackageName string
	TableName   string
}

type Interface struct {
	Name        string
	Methods     []Method
	PackageName string
}

type Method struct {
	Name       string
	Params     []Param
	Returns    []Return
	SQLComment string
}

type Param struct {
	Name string
	Type string
}

type Return struct {
	Type      string
	IsSlice   bool
	IsPointer bool
	BaseType  string
}

type Parser struct {
	fset       *token.FileSet
	structs    map[string][]*Struct
	interfaces map[string][]*Interface
	pkgName    string
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
				typeName := typeSpec.Name.Name
				if !token.IsExported(typeName) {
					continue
				}

				switch t := typeSpec.Type.(type) {
				case *ast.StructType:
					s := &Struct{
						Name:        typeName,
						PackageName: v.p.pkgName,
						TableName:   util.ToSnakeCase(typeName),
					}
					v.p.extractFields(t, s)
					if len(s.Fields) > 0 {
						if v.p.pkgName == "" {
							v.p.pkgName = "main"
						}
						v.p.structs[v.p.pkgName] = append(v.p.structs[v.p.pkgName], s)
					}
				case *ast.InterfaceType:
					iface := &Interface{
						Name:        typeName,
						PackageName: v.p.pkgName,
					}
					v.p.extractMethods(t, iface)
					if len(iface.Methods) > 0 {
						if v.p.pkgName == "" {
							v.p.pkgName = "main"
						}
						v.p.interfaces[v.p.pkgName] = append(v.p.interfaces[v.p.pkgName], iface)
					}
				}
			}
		}
	}
	return v
}

func New() *Parser {
	return &Parser{
		fset:       token.NewFileSet(),
		structs:    make(map[string][]*Struct),
		interfaces: make(map[string][]*Interface),
	}
}

func (p *Parser) ParseFiles(paths []string) error {
	p.structs = make(map[string][]*Struct)
	p.interfaces = make(map[string][]*Interface)
	p.pkgName = ""
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

	// Store the package name from the AST file
	if f.Name != nil {
		p.pkgName = f.Name.Name
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
			f.Type = exprToFullString(field.Type)
		}

		if f.Type == "" {
			f.Type = "interface{}"
		}

		s.Fields = append(s.Fields, f)
	}
}

func (p *Parser) extractMethods(ifaceType *ast.InterfaceType, iface *Interface) {
	if ifaceType.Methods == nil {
		return
	}

	for _, field := range ifaceType.Methods.List {
		funcType, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		if len(field.Names) == 0 {
			continue
		}

		m := Method{
			Name: field.Names[0].Name,
		}

		// Parse //sql:"..." from Doc comments
		if field.Doc != nil {
			for _, comment := range field.Doc.List {
				text := strings.TrimSpace(comment.Text)
				if strings.HasPrefix(text, "//sql:") {
					sqlVal := strings.TrimPrefix(text, "//sql:")
					// Use strconv.Unquote so that escape sequences like \"
					// inside the SQL string are properly resolved.
					if unquoted, err := strconv.Unquote(sqlVal); err == nil {
						sqlVal = unquoted
					} else {
						sqlVal = strings.Trim(sqlVal, `"`)
					}
					m.SQLComment = sqlVal
				}
			}
		}

		// Extract parameters (skip first ctx context.Context)
		if funcType.Params != nil {
			skipFirst := true
			for _, param := range funcType.Params.List {
				typeStr := exprToFullString(param.Type)
				if skipFirst && typeStr == "context.Context" {
					skipFirst = false
					continue
				}
				skipFirst = false

				if len(param.Names) == 0 {
					m.Params = append(m.Params, Param{Type: typeStr})
				} else {
					for _, name := range param.Names {
						m.Params = append(m.Params, Param{
							Name: name.Name,
							Type: typeStr,
						})
					}
				}
			}
		}

		// Extract return types
		if funcType.Results != nil {
			for _, result := range funcType.Results.List {
				typeStr := exprToFullString(result.Type)
				ret := Return{
					Type: typeStr,
				}

				// Analyze the type structure
				ret.IsSlice = strings.HasPrefix(typeStr, "[]*") || strings.HasPrefix(typeStr, "[]")
				ret.IsPointer = strings.HasPrefix(typeStr, "*")

				// Extract base type
				base := typeStr
				if strings.HasPrefix(base, "[]*") {
					base = strings.TrimPrefix(base, "[]*")
					ret.IsSlice = true
					ret.IsPointer = true
				} else if strings.HasPrefix(base, "[]") {
					base = strings.TrimPrefix(base, "[]")
				} else if strings.HasPrefix(base, "*") {
					base = strings.TrimPrefix(base, "*")
				}
				ret.BaseType = base

				if len(result.Names) == 0 {
					m.Returns = append(m.Returns, ret)
				} else {
					for range result.Names {
						m.Returns = append(m.Returns, ret)
					}
				}
			}
		}

		iface.Methods = append(iface.Methods, m)
	}
}

func (p *Parser) GetStructs() map[string][]*Struct {
	return p.structs
}

func (p *Parser) GetInterfaces() map[string][]*Interface {
	return p.interfaces
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

// exprToFullString converts an AST expression to a full type string,
// preserving package qualifiers (e.g., "context.Context", "sql.DB").
func exprToFullString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return exprToFullString(v.X) + "." + v.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToFullString(v.X)
	case *ast.ArrayType:
		return "[]" + exprToFullString(v.Elt)
	case *ast.MapType:
		return "map[" + exprToFullString(v.Key) + "]" + exprToFullString(v.Value)
	case *ast.ChanType:
		return "chan"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "..." + exprToFullString(v.Elt)
	default:
		return "any"
	}
}

func parseDBTag(tag string) string {
	// Parse struct tags which have format: key1:"value1" key2:"value2"
	// Find the db: key and return its value
	parts := parseTagParts(tag)
	for _, part := range parts {
		if strings.HasPrefix(part, "db:") {
			value := strings.TrimPrefix(part, "db:")
			value = strings.Trim(value, `"`)
			return value
		}
	}
	return ""
}

// parseTagParts splits a struct tag into key:value parts
// e.g., `db:"value" json:"field"` -> ["db:\"value\"", "json:\"field\""]
func parseTagParts(tag string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(tag); i++ {
		c := tag[i]
		switch c {
		case '"':
			inQuotes = !inQuotes
			current.WriteByte(c)
		case ' ':
			if inQuotes {
				current.WriteByte(c)
			} else if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}
