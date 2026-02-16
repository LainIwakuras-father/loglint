package main

import (
	"github.com/LainIwakuras-father/selectel-test-lint/addcheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(addcheck.Analyzer)
}
