package main

import (
	"github.com/LainIwakuras-father/loglint/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.Analyzer}
}

var AnalyzerPlugin analyzerPlugin
