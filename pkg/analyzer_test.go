package pkg

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	// analysistest.TestData() возвращает путь к директории testdata
	testdata := analysistest.TestData()

	// Запускаем анализатор на всех пакетах внутри testdata/src
	// Второй аргумент — путь к testdata, третий — анализатор, четвёртый — имена пакетов для проверки
	analysistest.Run(t, testdata, Analyzer, "basic")
}
