package rules

import (
	"go/token"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func CheckNoSpecial(pass *analysis.Pass, pos token.Pos, str string) {
	for _, r := range str {
		if !unicode.Is(unicode.Latin, r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			pass.Reportf(pos, "log message contains special characters or emojis")
			return
		}
	}
}
