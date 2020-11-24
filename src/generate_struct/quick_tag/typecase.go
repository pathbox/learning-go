package quicktag

import (
	"reflect"
	"unsafe"
)

type emptyInterface struct {
	pt unsafe.Pointer // 类型指针
	pv unsafe.Pointer // 值指针
}

func PointerOfType(t reflect.Type) unsafe.Pointer {
	p := *(*emptyInterface)(unsafe.Pointer(&t))
	return p.pv
}

func TypeCast(src interface{}, dstType reflect.Type) (dst interface{}) {
	srcType := reflect.TypeOf(src)
	eface := *(*emptyInterface)(unsafe.Pointer(&src)) // 指针的类型转换
	if srcType.Kind() == reflect.Ptr {
		eface.pt = PointerOfType(reflect.PtrTo(dstType))
	} else {
		eface.pt = PointerOfType(dstType)
	}
	dst = *(*interface{})(unsafe.Pointer(&eface))
	return
}
