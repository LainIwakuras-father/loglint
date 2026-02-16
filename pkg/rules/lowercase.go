package rules

import (
	"go/token"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func CheckLowercase(pass *analysis.Pass, pos token.Pos, str string) {
	if len(str) == 0 {
		return
	}
	first := str[0]
	if unicode.IsLetter(rune(first)) && !unicode.IsLower(rune(first)) {
		pass.Reportf(pos, "log message should start with a lowercase letter")
	}
}
