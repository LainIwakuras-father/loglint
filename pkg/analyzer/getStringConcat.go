package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// extractStringFromExpr пытается извлечь строковое значение из выражения expr.
// Обрабатывает:
//   - строковые литералы,
//   - именованные константы строкового типа,
//   - конкатенацию таких выражений через оператор + (если все части константны).
//
// Возвращает извлечённую строку и true в случае успеха, иначе false.
func extractStringFromExpr(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	// Сначала пробуем получить значение через TypesInfo (работает для любых констант)
	if str, ok := extractStringConstant(pass, expr); ok {
		return str, true
	}

	// Если не константа, проверяем, не конкатенация ли это
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok || bin.Op != token.ADD {
		return "", false
	}

	// Убеждаемся, что оба операнда — строки (чтобы не путать со сложением чисел)
	if !isStringType(pass, bin.X) || !isStringType(pass, bin.Y) {
		return "", false
	}

	// Рекурсивно извлекаем левую и правую части
	left, okLeft := extractStringFromExpr(pass, bin.X)
	right, okRight := extractStringFromExpr(pass, bin.Y)
	if !okLeft || !okRight {
		return "", false
	}

	return left + right, true
}

// extractStringConstant — вспомогательная функция, получающая строковое значение
// из константного выражения через TypesInfo.
func extractStringConstant(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok || tv.Value == nil || tv.Value.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(tv.Value), true
}

// isStringType проверяет, имеет ли выражение тип string.
func isStringType(pass *analysis.Pass, expr ast.Expr) bool {
	typ := pass.TypesInfo.TypeOf(expr)
	if typ == nil {
		return false
	}
	return types.Identical(typ, types.Typ[types.String])
}

func extract(pass *analysis.Pass, msg ast.Expr) {
	var isStatic func(msg ast.Expr) bool
	isStatic = func(msg ast.Expr) bool {
		switch msg := msg.(type) {
		case *ast.BasicLit: // e.g. slog.Info("msg")
			return msg.Kind == token.STRING
		case *ast.Ident: // e.g. slog.Info(constMsg)
			_, isConst := pass.TypesInfo.ObjectOf(msg).(*types.Const)
			return isConst
		case *ast.BinaryExpr: // e.g. slog.Info("x" + "y")
			if msg.Op != token.ADD {
				panic("unreachable") // Only "+" can be applied to strings.
			}
			return isStatic(msg.X) && isStatic(msg.Y)
		default:
			return false
		}
	}

	if !isStatic(msg) {
		pass.ReportRangef(msg, "message should be a string literal or a constant")
	}
}

func extractMessage(pass *analysis.Pass, call *ast.CallExpr) (string, token.Pos, bool) {
	if len(call.Args) == 0 {
		return "", token.NoPos, false
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", token.NoPos, false
	}
	// Убираем кавычки
	str := lit.Value
	if len(str) >= 2 {
		str = str[1 : len(str)-1]
	}
	return str, lit.Pos(), true
}
