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

func changeFuncAssignmentToSetCall(root ast.Node) int {
	count := 0
	astutil.Apply(root, nil, func(c *astutil.Cursor) bool {
		a, ok := c.Node().(*ast.AssignStmt)
		if ok {
			s, ok := a.Lhs[0].(*ast.SelectorExpr)
			if ok {
				if s.Sel.Name == "GetFunc" || s.Sel.Name == "FetchFunc" {
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
					count++
				}
			}
		}
		return true
	})
	return count
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal(`
rewrite <go_source_file> <destination_file>
`)
	}
	source := os.Args[1]
	destination := os.Args[2]
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	count := changeFuncAssignmentToSetCall(node)
	log.Printf("Changed %d mock function assignments to Set calls\n", count)

	f, err := os.Create(destination)
	defer f.Close()
	if err := printer.Fprint(f, fset, node); err != nil {
		log.Fatal(err)
	}

}
