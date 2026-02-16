package analyzer

import (
	"go/ast"
	"go/constant"

	"golang.org/x/tools/go/analysis"
)

func extractStringExpendKind(pass *analysis.Pass, call *ast.CallExpr) (string, bool) {
	// Пробуем получить тип выражения
	arg := call.Args[0]

	tv, ok := pass.TypesInfo.Types[arg]
	if !ok {
		return "", false
	}
	// Проверяем, что это константа и её тип — string
	if tv.Value == nil || tv.Value.Kind() != constant.String {
		return "", false
	}
	// Извлекаем строковое значение
	return constant.StringVal(tv.Value), true
}
