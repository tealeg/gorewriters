package main

import (
	"go/parser"
	"go/printer"

	"go/token"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestIdentifyMocks(t *testing.T) {
	c := qt.New(t)
	fset := token.NewFileSet()
	source := `package foo

import "foo/mock"

func TestGet(t *testing.T) {
	fooMock := mock.NewFooMock(t)
	barMock := newNotRealMock(t)
	var bazMock BazMock
	bazMock = mock.NewBazMock(t)

}
`
	node, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	c.Assert(err, qt.IsNil)
	mocks := identifyMocks(node)
	expected := &mockMap{"fooMock": true, "bazMock": true}
	c.Assert(mocks, qt.DeepEquals, expected)
}

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
	changeFuncAssignmentToSetCall(node, &mockMap{"tMock": true})

	var out strings.Builder
	if err := printer.Fprint(&out, fset, node); err != nil {
		c.Fatal(err)
	}
	c.Assert(out.String(), qt.Equals, expected)
}
