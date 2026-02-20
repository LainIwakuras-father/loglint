package rules

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// hasSensitiveKeywords проверяет, содержит ли строка ключевые слова (регистронезависимо).
func hasSensitiveKeywords(s string) bool {
	keywords := []string{
		"password", "passwd", "pwd",
		"token", "api_key", "apikey", "secret", "key",
	}
	lower := strings.ToLower(s)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// isSensitiveIdentifier проверяет, является ли имя идентификатора подозрительным.
func isSensitiveIdentifier(name string) bool {
	sensitiveNames := []string{
		"password", "passwd", "pwd",
		"token", "apiKey", "api_key", "secret", "key",
	}
	for _, sn := range sensitiveNames {
		if strings.EqualFold(name, sn) {
			return true
		}
	}
	return false
}

// containsSensitiveConcat рекурсивно обходит выражение и ищет опасные элементы.
func containsSensitiveConcat(pass *analysis.Pass, expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return false
		}
		//if !isStringType(pass, e.X) || !isStringType(pass, e.Y) {
		//	return false
		//}
		return containsSensitiveConcat(pass, e.X) || containsSensitiveConcat(pass, e.Y)

	case *ast.BasicLit:
		if e.Kind == token.STRING {
			str, err := strconv.Unquote(e.Value)
			if err != nil {
				return false
			}
			return hasSensitiveKeywords(str)
		}
		return false

	case *ast.Ident:
		return isSensitiveIdentifier(e.Name)

	default:
		return false
	}
}

// checkSensitive анализирует сообщение на наличие чувствительных данных.
func checkSensitive(pass *analysis.Pass, expr ast.Expr, msg string, ok bool) {
	if ok {
		// Удалось извлечь константную строку
		if hasSensitiveKeywords(msg) {
			pass.Reportf(expr.Pos(), "log message contains sensitive data (password, token, etc.)")
		}
	} else {
		// Не удалось извлечь строку — анализируем AST
		if containsSensitiveConcat(pass, expr) {
			pass.Reportf(expr.Pos(), "log message may contain sensitive data (concatenation with sensitive variable or literal)")
		}
	}
}
