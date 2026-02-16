// Package addcheck defines an Analyzer that reports time package expressions that
// can be simplified
package analyzer

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"

	r "github.com/LainIwakuras-father/loglint/pkg/rules"
	"golang.org/x/tools/go/analysis"
)

//var Analyzer = &analysis.Analyzer{
//	Name: "addlint",
//	Doc:  "reports integer additions",
//	Run:  run,
//}

func NewAnalyzer() (*analysis.Analyzer, error) {
	return &analysis.Analyzer{
		Name: "loglint",
		Doc:  "reports integer additions",
		URL:  "https://github.com/LainIwakuras-father/loglint",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return run(pass)
		},
	}, nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	// find in file all call
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			// veryfy call
			if isLogCall(pass, call) {
				return true
			}
			// extracting message from log
			msg, ok := extractStringExpendKind(pass, call)
			if ok {
				// pass
				r.CheckLowercase(pass, call.Pos(), msg)
				r.CheckEnglish(pass, call.Pos(), msg)
				r.CheckNoSpecial(pass, call.Pos(), msg)
				return true
			}
			return true // compliance check
		})
	}

	return nil, nil
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
