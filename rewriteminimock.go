package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "../foo.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	astutil.Apply(node, nil, func(c *astutil.Cursor) bool {
		a, ok := c.Node().(*ast.AssignStmt)
		if ok {
			s, ok := a.Lhs[0].(*ast.SelectorExpr)
			if ok {
				if s.Sel.Name == "GetFunc" {
					fn, ok := a.Rhs[0].(*ast.FuncLit)
					if !ok {
						log.Fatalf("RHS is not a FuncLit, it's a: %T %+v\n", a.Rhs[0], a.Rhs[0])
						return false
					}

					newNode := &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   s.X,
								Sel: ast.NewIdent("Set"),
							},
							Args: []ast.Expr{
								fn,
							},
						},
					}
					c.Replace(newNode)
				}
			}
		}
		return true
	})

	f, err := os.Create("new.go")
	defer f.Close()
	if err := printer.Fprint(f, fset, node); err != nil {
		log.Fatal(err)
	}

}
