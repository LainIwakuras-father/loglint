package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func isLogCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	// find call with help SelectorExpr
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// find type in file
	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil {
		return false
	}
	funcObj, ok := obj.(*types.Func)
	if !ok {
		return false
	}
	// find where func??
	pkg := funcObj.Pkg()
	if pkg == nil {
		return false
	}
	// Verify path package
	switch pkg.Path() {
	case "log/slog", "go.uber.org/zap":
		// Проверяем имя метода
		switch funcObj.Name() {
		case "Info", "Error", "Debug", "Warn", "Fatal", "Panic":
			return true
		// Для zap есть сахарные методы
		case "Infow", "Errorw", "Debugw", "Warnw", "Fatalw", "Panicw":
			return true
		}
	}

	return false
}
