package analyzer

import (
	"go/ast"
	"go/types"

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
			if !isHTTPResponseCall(pass, call) {
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

func isHTTPResponseCall(pass *analysis.Pass, call *ast.CallExpr) bool {

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	obj := pass.TypesInfo.ObjectOf(selector.Sel)
	if obj == nil || obj.Pkg() == nil {
		return false
	}
	if obj.Pkg().Path() != "net/http" {
		return false
	}
	fn, ok := obj.(*types.Func)
	if !ok {
		return false
	}
	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		return false
	}
	if recv := sig.Recv(); recv != nil {
		recvType := recv.Type()

		if ptr, ok := recvType.(*types.Pointer); ok {
			recvType = ptr.Elem()
		}

		named, ok := recvType.(*types.Named)
		if !ok {
			return false
		}

		if named.Obj().Pkg() == nil ||
			named.Obj().Pkg().Path() != "net/http" ||
			named.Obj().Name() != "Client" {
			return false
		}

		switch fn.Name() {
		case "Do", "Get":
			return true
		default:
			return false
		}
	}
	switch selector.Sel.Name {
	case "Get", "Post", "PostForm", "Head":
		return true
	default:
		return false
	}

}
