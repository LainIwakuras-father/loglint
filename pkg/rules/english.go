package rules

import (
	"go/token"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// checkEnglish проверяет, что сообщение содержит только латинские буквы, цифры и пробелы.
func CheckEnglish(pass *analysis.Pass, pos token.Pos, msg string) {
	for _, r := range msg {
		if !unicode.Is(unicode.Latin, r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			pass.Reportf(pos, "log message must contain only English letters, digits, and spaces")
			return // можно выйти после первого нарушения, или собирать все, но для простоты достаточно одного
		}
	}
}
