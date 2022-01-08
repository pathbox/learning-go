package decorator

import "fmt"

type Booker interface {
	Reading()
}

type Book struct{}

func (book Book) Reading() {
	fmt.Println("reading")
}

type Underliner interface {
	Booker
	Underline()
}

type NotesTaker interface {
	Booker
	TakeNotes()
}

type ConcreteUnderline struct {
	Booker Booker // 包裹了Booker
}

// 提供读书的方法，包装了Booker接口
func (underline ConcreteUnderline) Reading() {
	underline.Booker.Reading()
}

/*
在运行期间动态地给对象添加额外的职责，比子类更灵活
若不想修改原来的接口，则可以使用装饰者
允许向一个现有的对象添加新的功能，同时又不改变其结构。这种类型的设计模式属于结构型模式，它是作为现有的类的一个包装。

这种模式创建了一个装饰类，用来包装原有的类，并在保持类方法签名完整性的前提下，提供了额外的功能
*/
