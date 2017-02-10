package main

import(  
"fmt"  
"flag"  
"strconv"  
)

type CustomStruct struct {  
    FlagValue int
}

func (cF *CustomStruct) String() string {  
    return strconv.Itoa(cF.FlagValue)
}

func (cF *CustomStruct) Set(s string) error {  
    cF.FlagValue, _ = strconv.Atoi(s)

    return nil
}

func (cF *CustomStruct) Get() int {  
    return cF.FlagValue
}

func (cF *CustomStruct) IsSet() bool {  
    if cF.FlagValue == -1 {
        return false
    }

    return true
}

func main() {  
    limitFlag := CustomStruct{-1}

    flag.Var(&limitFlag, "limit", "Limits output")

    flag.Parse()

    if limitFlag.IsSet() {
        fmt.Printf("\nLimit: %d\n\n", limitFlag.Get())
    } else {
        fmt.Printf("\nLimit flag not included.\n\n")
    }
}