package doc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/printer"
	"go/token"
	"strings"

	"gondola/internal/astutil"
	"gondola/template"
)

const (
	constPrefix = "const-"
	varPrefix   = "var-"
)

func ConstId(name string) string {
	return constPrefix + name
}

func VarId(name string) string {
	return varPrefix + name
}

func FuncId(name string) string {
	return "func-" + name
}

func TypeId(name string) string {
	return "type-" + name
}

func MethodId(typ string, name string) string {
	return "type-" + typ + "-method-" + name
}

func trim(s string, t string) string {
	return strings.Trim(s, t)
}

func funcId(fn *doc.Func) string {
	if fn.Recv != "" {
		recv := fn.Recv
		if recv[0] == '*' {
			recv = recv[1:]
		}
		return MethodId(recv, fn.Name)
	}
	return FuncId(fn.Name)
}

func FuncReceiver(fn *ast.FuncDecl) string {
	if fn.Recv != nil {
		recv := astutil.Ident(fn.Recv.List[0].Type)
		if recv[0] == '*' {
			recv = recv[1:]
		}
		return recv
	}
	return ""
}

func funcName(fn *ast.FuncDecl) string {
	if recv := FuncReceiver(fn); recv != "" {
		return "(" + recv + ") " + fn.Name.Name
	}
	return fn.Name.Name
}

func typeId(typ *doc.Type) string {
	return TypeId(typ.Name)
}

func documentedColor(val float64) string {
	if val < 25 {
		return "red"
	}
	if val < 50 {
		return "yello"
	}
	if val < 75 {
		return "lightblue"
	}
	return "green"
}

func issuesColor(val int) string {
	if val > 20 {
		return "red"
	}
	if val > 10 {
		return "yellow"
	}
	return "green"
}

func complexityColor(val interface{}) (string, error) {
	var v float64
	switch va := val.(type) {
	case int:
		v = float64(va)
	case float64:
		v = va
	default:
		return "", fmt.Errorf("invalid complexity type %T", val)
	}
	if v < 5 {
		return "green", nil
	}
	if v < 10 {
		return "yellow", nil
	}
	return "red", nil
}

func bootstrapColor(c string) string {
	switch c {
	case "green":
		return "success"
	case "lightblue":
		return "info"
	case "yellow":
		return "warning"
	case "red":
		return "danger"
	}
	return ""
}

func funcListName(fn *ast.FuncDecl) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), fn)
	return strings.TrimPrefix(buf.String(), "func ")
}

func init() {
	template.AddFuncs([]*template.Func{
		{Name: "func_id", Fn: funcId},
		{Name: "func_name", Fn: funcName},
		{Name: "type_id", Fn: typeId},
		{Name: "issues_color", Fn: issuesColor},
		{Name: "complexity_color", Fn: complexityColor},
		{Name: "bootstrap_color", Fn: bootstrapColor},
		{Name: "func_list_name", Fn: funcListName},
	})
}
