package main

import "C"

/*
#include <stdio.h>

static void SayHello(const char* s) {
	puts(s);
}
*/
func main() {
	C.SayHello(C.CString("Hello World\n"))
}
