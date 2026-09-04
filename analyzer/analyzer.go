package analyzer

import (
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "cleanupcheck",
	Doc:  "check for resources that are acquired but not properly cleaned up",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	checkHTTPResponseBody(pass)

	return nil, nil
}
