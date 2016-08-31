package main

/*
#include <stdio.h>
int datspecialnumber() {
    return 2015;
}
*/

import "C"
import "fmt"

func main() {
    fmt.Println(C.datspecialnumber())
}
