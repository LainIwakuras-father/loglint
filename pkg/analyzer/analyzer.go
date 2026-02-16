package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "checks log messages for style and security issues",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		checkLogCall(pass, call)
	})

	return nil, nil
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	pkgName := pkgIdent.Name
	methodName := sel.Sel.Name

	if !isLoggingMethod(pkgName, methodName) {
		return
	}

	msgArg := getMessageArg(call)
	if msgArg == nil {
		return
	}

	checkLowercase(pass, msgArg)
	checkEnglish(pass, msgArg)
	checkNoSpecial(pass, msgArg)
	checkSensitive(pass, msgArg)
}

func isLoggingMethod(pkgName, methodName string) bool {
	loggingPackages := map[string][]string{
		"log":  {"Info", "Error", "Debug", "Warn", "Print", "Printf"},
		"slog": {"Info", "Error", "Debug", "Warn"},
		"zap":  {"Info", "Error", "Debug", "Warn"},
	}

	methods, exists := loggingPackages[pkgName]
	if !exists {
		return false
	}

	for _, m := range methods {
		if methodName == m {
			return true
		}
	}
	return false
}

func getMessageArg(call *ast.CallExpr) ast.Expr {
	if len(call.Args) == 0 {
		return nil
	}
	return call.Args[0]
}

func getStringValue(expr ast.Expr) (string, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			value, err := strconv.Unquote(v.Value)
			if err != nil {
				return "", false
			}
			return value, true
		}
	case *ast.BinaryExpr:
		if v.Op == token.ADD {
			left, ok := getStringValue(v.X)
			if !ok {
				return "", false
			}
			right, ok := getStringValue(v.Y)
			if !ok {
				return "", false
			}
			return left + right, true
		}
	}
	return "", false
}
