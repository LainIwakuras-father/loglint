package analyzer

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("myanalyzer", New)
}

type settings struct {
	// сюда можно принимать настройки из .golangci.yml
	ExtraMessage string `json:"extra-message"`
}

func New(conf any) (register.LinterPlugin, error) {
	//	var s settings
	// парсинг конфига (опционально)
	// ...

	return register.LinterPlugin{
		Analyzers:   []*analysis.Analyzer{Analyzer},
		Name:        "myanalyzer",
		Description: "Запрещает голые panic",
		// Config:         s,  // если нужны настройки
	}, nil
}
