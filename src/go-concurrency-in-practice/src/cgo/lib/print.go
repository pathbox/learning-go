package lib

/*
#include <stdio.h>
#include <stdlib.h>
void myprint(char* s) {
        printf("%s", s);
}
*/

import "C"

func Print(s string) {
	cs := C.CSting(s)
	defer C.free(unsage.Pointer(cs)) // 释放指针
	C.myprint(cs)
}
