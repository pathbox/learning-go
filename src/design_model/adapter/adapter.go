package adapter

<<<<<<< HEAD
import "fmt" // 适配器模式
=======
import "fmt"

type Volts220 struct{}

func (v Volts220) OutputPower() {
	fmt.Println("电源电压220V")
}

type Adaptee interface {
	OutputPower()
}

type Target interface {
	ConvertTo5V()
}

type Adapter struct {
	Adaptee // 原有的适配器 Adapter struct 实现了Target interface的方法
}

func (a Adapter) ConvertTo5V() {
	a.OutputPower()
	fmt.Println("转为5V电压")
}

>>>>>>> 3d64ae44130a20de20c1a1760f430acfeba757e5
