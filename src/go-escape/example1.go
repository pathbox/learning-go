package main

import "fmt"

type S struct { 
 M *int
}

func main() { 
 var x S
 var i int
 ref(&i, &x) // &i escapes to heap
 fmt.Println(x.M) // x.M escapes to heap
}

func ref(y *int, z *S) { 
 z.M = y
}