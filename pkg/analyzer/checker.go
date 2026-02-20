package analyzer

import (
	"go/ast"

	r "github.com/LainIwakuras-father/loglint/pkg/rules"
	"golang.org/x/tools/go/analysis"
)

func run(pass *analysis.Pass) (any, error) {
	// find in file all call
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			// veryfy call
			if !isLogCall(pass, call) {
				return true
			}
			// extracting message from log
			msg, ok := extractStringExpendKind(pass, call)
			if !ok {
				return true
			}
			// pass
			r.CheckLowercase(pass, call.Pos(), msg)
			// r.CheckEnglish(pass, call.Pos(), msg)
			r.CheckNoSpecial(pass, call.Pos(), msg)
			return true
		})
	}

	return nil, nil
}
