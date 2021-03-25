


指针类型：

*类型：普通指针，用于传递对象地址，不能进行指针运算。

unsafe.Pointer：通用指针类型，用于转换不同类型的指针，不能进行指针运算。

uintptr：用于指针运算，GC 不把 uintptr 当指针，uintptr 无法持有对象。uintptr 类型的目标会被回收。

　　unsafe.Pointer 可以和 普通指针 进行相互转换。
　　unsafe.Pointer 可以和 uintptr 进行相互转换。

　　也就是说 unsafe.Pointer 是桥梁，可以让任意类型的指针实现相互转换，也可以将任意类型的指针转换为 uintptr 进行指针运算。

转为 uintptr 后就能进行偏移计算，unsafe.Offsetof(s.d)的结果是 uintptr


结构体成员的内存分配是连续的，第一个成员的地址就是结构体的地址，相对于结构体的偏移量为 0。其它成员都可以通过偏移量来计算其地址。

每种类型都有它的大小和对齐值，可以通过 unsafe.Sizeof 获取其大小，通过 unsafe.Alignof 获取其对齐值，通过 unsafe.Offsetof 获取其偏移量。不过 unsafe.Alignof 获取到的对齐值只是该类型单独使用时的对齐值，不是作为结构体字段时与其它对象间的对齐值，这里用不上，所以需要用 unsafe.Offsetof 来获取字段的偏移量，进而确定其内存地址