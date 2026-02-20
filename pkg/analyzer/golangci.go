package analyzer

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

// ← Это основной файл плагина для golangci-lint
func init() {
	register.Plugin("loglint", New) // имя линтера должно совпадать с тем, что enable в .golangci.yml
	// или "loglint" — как хочешь, но一致но
}

type plugin struct{} // ← любое имя, часто просто plugin или linterPlugin

// Явная проверка реализации интерфейса (очень полезно)
var _ register.LinterPlugin = (*plugin)(nil)

func New(settings any) (register.LinterPlugin, error) {
	// settings — это map[string]any из .golangci.yml → linters-settings.custom.addlint
	// если настройки не нужны — просто игнорируем
	return &plugin{}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		Analyzer, // ← именно эта переменная!
	}, nil
}

func (p *plugin) GetLoadMode() string {
	// Самый распространённый и достаточный вариант для 90% линтеров
	return register.LoadModeTypesInfo
	// Если нужен только синтаксис (быстрее) → register.LoadModeSyntax
	// Если очень сложный анализ → register.LoadModeFull (медленно)
}
