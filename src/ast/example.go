package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"reflect"
)
// This example shows to a Go program can parse itself into an AST tree and
// to infer types of nodes (expressions) in this AST tree.

const (
  kPath = "./"
)

func main() {
  fset := token.NewFileSet()

  pkgs, e := parser.ParseDir(fset, kPath, nil, 0)
  if e != nil {
    log.Fatal(e)
    return
  }

  astf := make([]*ast.File, 0)
  for _, pkg := range pkgs{
    fmt.Printf("package %v\n", pkg.Name)
    for fn, f := range pkg.Files{
      fmt.Printf("file %v\n", fn)
      astf = append(astf, f)
    }
  }
  config := &types.Config{
    Error: func(e error){
      fmt.Println(e)
    },
    Importer: importer.Default(),
  }
  info := types.Info{
    Types: make(map[ast.Expr]types.TypeAndValue),
    Defs: make(map[*ast.Ident]types.Object),
    Uses: make(map[*ast.Ident]types.Object),
  }
  pkg, e := config.Check(kPath, fset, astf, &info)
  if e!= nil{
    fmt.Println(e)
  }
  fmt.Printf("types.Config.Check got %v\n", pkg.String())
  for _, f:= range astf{
    ast.Walk(&PrintASTVisitor{&info}, f)
  }
}

type PrintASTVisitor struct{
  info *types.Info
}

func (v *PrintASTVisitor) Visit(node ast.Node) ast.Visitor  {
  if node != nil {
    fmt.Printf("%s", reflect.TypeOf(node).String())
    switch node.(type) {
    case ast.Expr:
      t := v.info.TypeOf(node.(ast.Expr))
      if t != nil {
        fmt.Printf(": %s", t.String())
      }
    }
    fmt.Println()
  }
  return v
}
