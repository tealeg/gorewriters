package main

import (
	"go/parser"
	"go/printer"

	"go/token"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestChangeFuncAssignmentToSetCall(t *testing.T) {
	c := qt.New(t)
	fset := token.NewFileSet()
	source := `package foo

import "foo/mock"

func TestGet(t *testing.T) {
	tMock := mock.NewFooMock(t)
	tMock.GetFunc = func() bool {
		return true
	}
}
`
	expected := `package foo

import "foo/mock"

func TestGet(t *testing.T) {
	tMock := mock.NewFooMock(t)
	tMock.Set(func() bool {
		return true
	})
}
`

	node, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	c.Assert(err, qt.IsNil)
	changeFuncAssignmentToSetCall(node)

	var out strings.Builder
	if err := printer.Fprint(&out, fset, node); err != nil {
		c.Fatal(err)
	}
	c.Assert(out.String(), qt.Equals, expected)

}
