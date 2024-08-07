package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func GetExitAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "os_exit_checker",
		Doc:  "checks for calls to os.Exit in the main function of package main",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// first, find if this is the main() inside package main
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}
			// if this is the main() inside package main, look for os.Exit() calls inside
			if funcDecl.Name.Name == "main" && pass.Pkg.Name() == "main" {
				// inspect the main() function body for os.Exit() calls
				ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
					callExpr, ok := n.(*ast.CallExpr)
					if !ok {
						return true
					}
					selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
					if !ok {
						return true
					}
					exitPkgIdent, ok := selExpr.X.(*ast.Ident)
					if !ok {
						return true
					}
					if exitPkgIdent.Name == "os" && selExpr.Sel.Name == "Exit" {
						pass.Reportf(callExpr.Pos(), "calling os.Exit in the main function is forbidden in package main")
					}
					return true
				})
			}
			return true
		})
	}

	return nil, nil
}
