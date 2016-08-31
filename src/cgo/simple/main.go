package main

/*
#include <stdlib.h>
extern void my_puts(const char*);
*/
import "C"
import "unsafe"

func main() {
	p := C.CString("Golang is awsome")
	defer C.free(unsafe.Pointer(p))
	C.my_puts(p)
}
