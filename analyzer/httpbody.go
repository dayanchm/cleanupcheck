package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func checkHTTPResponseBody(pass *analysis.Pass) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return true
			}
			checkFunctionForHTTPBody(pass, fn)
			return false
		})
	}
}

func checkFunctionForHTTPBody(pass *analysis.Pass, fn *ast.FuncDecl) {
	responses := make(map[string]ast.Expr)

	//
	// First pass:
	// Find: resp, err := http.Get(...)

	ast.Inspect(fn.Body, func(node ast.Node) bool {

		assign, ok := node.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for i, rhs := range assign.Rhs {
			call, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}
			if !isHTTPResponseCall(call) {
				continue
			}
			if len(assign.Lhs) <= i {
				continue
			}
			ident, ok := assign.Lhs[i].(*ast.Ident)
			if !ok {
				continue
			}
			responses[ident.Name] = rhs
		}
		return true
	})

	if len(responses) == 0 {
		return
	}
	//
	// Second pass:
	// Find:
	//
	//     resp.Body.Close()
	//
	closed := make(map[string]bool)
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "Close" {
			return true
		}
		bodySelector, ok := selector.X.(*ast.SelectorExpr)
		if !ok || bodySelector.Sel.Name != "Body" {
			return true
		}
		responseIdent, ok := bodySelector.X.(*ast.Ident)
		if !ok {
			return true
		}
		if _, exists := responses[responseIdent.Name]; exists {
			closed[responseIdent.Name] = true
		}
		return true
	})
	for name, expr := range responses {
		if closed[name] {
			continue
		}
		pass.Reportf(
			expr.Pos(),
			"HTTP response %q is not closed; call defer %s.Body.Close()",
			name,
			name,
		)
	}

}

func isHTTPResponseCall(call *ast.CallExpr) bool {

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}
	if pkg.Name != "http" {
		return false
	}
	switch selector.Sel.Name {
	case "Get", "Post", "PostForm", "Head":
		return true
	default:
		return false
	}

}
