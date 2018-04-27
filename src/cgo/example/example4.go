// int32 => *C.char
var x := int32(9527)
var p *C.char = (*C.char) (unsafe.Pointer(uintptr(x)))

// *C.char => int32
var y *C.char
var q int32 = int32(uintptr(unsafe.Pointer(y)))

int32 <=> uintptr <=> unsafe.Pointer <=> *C.char

var p *X
var q *Y

q = (*Y)(unsafe.Pointer(p)) // *X => *Y
p = (*X)(unsafe.Pointer(q)) // *Y => *X

*X <=> unsafe.Pointer <=> *Y

var p []X
var q []Y

pHdr := (*reflect.SliceHeader)(unsafe.Pointer(&p))
qHdr := (*reflect.SliceHeader)(unsafe.Pointer(&q))

pHdr.Data = qHdr.Data
pHdr.Len = qHdr.Len * unsafe.Sizeof(q[0])
pHdr.Cap = qHdr.Cap * unsafe.Sizeof(q[0])


// []float64 强制类型转换为 []int
var a = []float64{4,5,6,7}
var b []int = ((*[1<<20]int))(unsafe.Pointer(&a[0]))[:len]

sort.Ints(b)