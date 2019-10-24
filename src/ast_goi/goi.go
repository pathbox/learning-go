package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"
)

const TRACE = false

type valueType int

const (
	i64 valueType = iota
	fn
	bl
)

type value struct {
	t   valueType
	val interface{}
}

type context map[string]*value

type ret struct {
	set  bool
	vals []value
}

func (r *ret) setVals(vals []value) {
	r.set = true
	r.vals = vals
}

func (r *ret) setVal(v value) {
	r.set = true
	r.vals = []value{v}
}

type fnType func(context, *ret, []value)

func (ctx context) copy() context {
	cpy := context{}
	for key, value := range ctx {
		cpy[key] = value
	}

	return cpy
}

func newContext() context {
	return map[string]*value{}
}

func interpretBinaryExpr(ctx context, r *ret, bexpr *ast.BinaryExpr) {
	var xr, yr ret
	interpretExpr(ctx, &xr, bexpr.X)
	x := xr.vals[0]
	interpretExpr(ctx, &yr, bexpr.Y)
	y := yr.vals[0]

	switch bexpr.Op {
	case token.ADD:
		r.setVal(value{i64, x.val.(int64) + y.val.(int64)})
		return
	case token.SUB:
		r.setVal(value{i64, x.val.(int64) - y.val.(int64)})
		return
	case token.MUL:
		r.setVal(value{i64, x.val.(int64) * y.val.(int64)})
		return
	case token.QUO:
		r.setVal(value{i64, x.val.(int64) / y.val.(int64)})
		return
	case token.LSS:
		r.setVal(value{i64, x.val.(int64) < y.val.(int64)})
		return
	case token.EQL:
		r.setVal(value{bl, x.val.(int64) == y.val.(int64)})
		return
	}

	panic("Unhandled op type")
}

func interpretCallExpr(ctx context, r *ret, ce *ast.CallExpr) {
	var fnr ret
	interpretExpr(ctx, &fnr, ce.Fun)
	fn := fnr.vals[0]

	vals := []value{}
	for _, arg := range ce.Args {
		var vr ret
		interpretExpr(ctx, &vr, arg)
		vals = append(vals, vr.vals[0])
	}

	fn.val.(fnType)(ctx, r, vals)
}

func interpretExpr(ctx context, r *ret, expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		interpretBinaryExpr(ctx, r, e)
		return
	case *ast.BasicLit:
		switch e.Kind {
		case token.INT:
			i, _ := strconv.ParseInt(e.Value, 10, 64)
			r.setVal(value{i64, i})
		}
		return
	case *ast.Ident:
		r.setVal(*ctx[e.Name])
		return
	case *ast.CallExpr:
		interpretCallExpr(ctx, r, e)
		return
	}

	fmt.Println(expr)
	panic("Unexpected expr type")
}

func interpretReturnStmt(ctx context, r *ret, s *ast.ReturnStmt) {
	var vals []value
	for _, expr := range s.Results {
		var r ret
		interpretExpr(ctx, &r, expr)
		vals = append(vals, r.vals[0])
	}

	r.setVals(vals)

	return
}

func interpretIfStmt(ctx context, r *ret, is *ast.IfStmt) {
	if is.Init != nil {
		interpretStmt(ctx, nil, is.Init)
	}

	var cr ret
	interpretExpr(ctx, &cr, is.Cond)
	c := cr.vals[0]
	if c.val.(bool) {
		interpretBlockStmt(ctx, r, is.Body)
		return
	}

	interpretStmt(ctx, r, is.Else)
}

func interpretForStmt(ctx context, r *ret, fs *ast.ForStmt) {
	var ir ret
	interpretStmt(ctx, &ir, fs.Init)

	for {
		var cr ret
		interpretExpr(ctx, &cr, fs.Cond)

		if !cr.vals[0].val.(bool) {
			break
		}

		interpretStmt(ctx, r, fs.Body)
		interpretStmt(ctx, nil, fs.Post)
	}
}

func interpretAssignStmt(ctx context, r *ret, as *ast.AssignStmt) {
	for i, lhs := range as.Lhs {
		rhs := as.Rhs[i]

		switch l := lhs.(type) {
		case *ast.Ident:
			// In other cases should be interpreted _after_ lhs
			var rr ret
			interpretExpr(ctx, &rr, rhs)

			if as.Tok == token.DEFINE {
				ctx[l.Name] = &rr.vals[0]
			} else {
				*ctx[l.Name] = rr.vals[0]
			}
		default:
			fmt.Println(lhs, rhs)
			panic("Unsupported assign")
		}
	}
}

func interpretStmt(ctx context, r *ret, stmt ast.Stmt) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.ReturnStmt:
		interpretReturnStmt(ctx, r, s)
		return
	case *ast.IfStmt:
		interpretIfStmt(ctx, r, s)
		return
	case *ast.BlockStmt:
		interpretBlockStmt(ctx, r, s)
		return
	case *ast.ForStmt:
		interpretForStmt(ctx, r, s)
		return
	case *ast.AssignStmt:
		interpretAssignStmt(ctx, r, s)
		return
	case *ast.ExprStmt:
		interpretExpr(ctx, r, s.X)
		return
	case *ast.IncDecStmt:
		switch e := s.X.(type) {
		case *ast.Ident:
			v := ctx[e.Name]
			if s.Tok == token.INC {
				v.val = v.val.(int64) + 1
			} else {
				v.val = v.val.(int64) - 1
			}
			return
		default:
			panic("Unsupported incdec type")
		}
	}

	fmt.Printf("%+v\n", stmt)
	panic("Unhandled stmt type")
}

func interpretBlockStmt(ctx context, r *ret, body *ast.BlockStmt) {
	for _, stmt := range body.List {
		interpretStmt(ctx, r, stmt)
		if r.set {
			return
		}
	}
}

func interpretFuncDecl(ctx context, r *ret, fd *ast.FuncDecl) {
	ctx[fd.Name.String()] = &value{
		fn,
		fnType(func(ctx context, r *ret, args []value) {
			if TRACE {
				fmt.Printf("TRACE: %s\n", fd.Name.String())
			}
			childCtx := ctx.copy()
			for i, arg := range args {
				argName := fd.Type.Params.List[i].Names[0].Name
				childCtx[argName] = &arg
			}
			interpretBlockStmt(childCtx, r, fd.Body)
		}),
	}
}

func interpret(f *ast.File) {
	ctx := newContext()
	ctx["println"] = &value{
		fn,
		fnType(func(ctx context, r *ret, args []value) {
			var values []interface{}
			for _, arg := range args {
				values = append(values, arg.val)
			}

			fmt.Println(values...)
		}),
	}

	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok {
			interpretFuncDecl(ctx, nil, fd)
		}
	}

	var r ret
	(*ctx["main"]).val.(fnType)(ctx, &r, []value{})
}

func main() {
	src, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	interpret(f)
}
