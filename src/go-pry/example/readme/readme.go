package main

import "github.com/d4l3k/go-pry/pry"

func main() {
	a := 1
	pry.Apply(&pry.Scope{Vals:map[string]interface{}{ "main": main, "a": a, "pry": pry.Package{Name: "pry", Functions: map[string]interface{}{"Fuzz": pry.Fuzz,"Package": pry.Type(pry.Package{}),"Pry": pry.Pry,"Apply": pry.Apply,"JSImporter": pry.Type(pry.JSImporter{}),"Type": pry.Type,"Highlight": pry.Highlight,"InterpretError": pry.Type(pry.InterpretError{}),"Append": pry.Append,"Make": pry.Make,"Close": pry.Close,"Len": pry.Len,"DeAssign": pry.DeAssign,"ComputeBinaryOp": pry.ComputeBinaryOp,"ErrChanRecvFailed": pry.ErrChanRecvFailed,"ErrChanRecvInSelect": pry.ErrChanRecvInSelect,"ErrDivisionByZero": pry.ErrDivisionByZero,"ErrChanSendFailed": pry.ErrChanSendFailed,"ErrBranchContinue": pry.ErrBranchContinue,"Scope": pry.Type(pry.Scope{}),"ValuesToInterfaces": pry.ValuesToInterfaces,"ErrBranchBreak": pry.ErrBranchBreak,"Defer": pry.Type(pry.Defer{}),"NewScope": pry.NewScope,"Func": pry.Type(pry.Func{}),"StringToType": pry.StringToType,}}, }})
	_ = a
}
