package main

import (
	"./fsm"
	"fmt"
)

var (
	Poweroff        = fsm.FSMState("关闭")
	FirstGear       = fsm.FSMState("1档")
	SecondGear      = fsm.FSMState("2档")
	ThirdGear       = fsm.FSMState("3档")
	PowerOffEvent   = fsm.FSMEvent("按下关闭按钮")
	FirstGearEvent  = fsm.FSMEvent("按下1档按钮")
	SecondGearEvent = fsm.FSMEvent("按下2档按钮")
	ThirdGearEvent  = fsm.FSMEvent("按下3档按钮")
	PowerOffHandler = fsm.FSMHandler(func() fsm.FSMState {
		fmt.Println("电风扇已关闭")
		return Poweroff
	})
	FirstGearHandler = fsm.FSMHandler(func() fsm.FSMState {
		fmt.Println("电风扇开启1档，微风徐来！")
		return FirstGear
	})
	SecondGearHandler = fsm.FSMHandler(func() fsm.FSMState {
		fmt.Println("电风扇开启2档，凉飕飕！")
		return SecondGear
	})
	ThirdGearHandler = fsm.FSMHandler(func() fsm.FSMState {
		fmt.Println("电风扇开启3档，发型被吹乱了！")
		return ThirdGear
	})
)

// 电风扇
type ElectricFan struct {
	*fsm.FSM
}

func NewElectricFan(initState fsm.FSMState) *ElectricFan {
	return &ElectricFan{
		FSM: fsm.NewFSM(initState),
	}
}

func main() {
	efan := NewElectricFan(Poweroff) // 初始状态是关闭状态

	// 关闭状态
	// 把所有可能的类别事件 存储到map中
	efan.AddHandler(Poweroff, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(Poweroff, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(Poweroff, SecondGearEvent, SecondGearHandler)
	efan.AddHandler(Poweroff, ThirdGearEvent, ThirdGearHandler)
	// 1档状态
	efan.AddHandler(FirstGear, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(FirstGear, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(FirstGear, SecondGearEvent, SecondGearHandler)
	efan.AddHandler(FirstGear, ThirdGearEvent, ThirdGearHandler)
	// 2档状态
	efan.AddHandler(SecondGear, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(SecondGear, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(SecondGear, SecondGearEvent, SecondGearHandler)
	efan.AddHandler(SecondGear, ThirdGearEvent, ThirdGearHandler)
	// 3档状态
	efan.AddHandler(ThirdGear, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(ThirdGear, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(ThirdGear, SecondGearEvent, SecondGearHandler)
	efan.AddHandler(ThirdGear, ThirdGearEvent, ThirdGearHandler)

	// 开始测试状态变化
	efan.Call(ThirdGearEvent)  // 按下3档按钮
	efan.Call(FirstGearEvent)  // 按下1档按钮
	efan.Call(PowerOffEvent)   // 按下关闭按钮
	efan.Call(SecondGearEvent) // 按下2档按钮
	efan.Call(PowerOffEvent)   // 按下关闭按钮
}

// 参考链接 http://www.jianshu.com/p/37281543f506
