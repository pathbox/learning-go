package main

import (
	"fmt"
	"math"
)

type areaError struct {
	err    string
	radius float64
	EType  int
}

// 实现errors的这个接口 则areaError struct 就是error类型
func (e *areaError) Error() string {
	return e.err
}

func circleArea(radius float64) (float64, error) {
	if radius < 0 {
		return 0, &areaError{"radius is negative", radius, 1}
	}
	return math.Pi * radius * radius, nil
}

func main() {
	radius := -20.0
	area, err := circleArea(radius)
	if err != nil {
		if err1, ok := err.(*areaError); ok { // 用接口方式转换，判断，得到的err1就是areaError 能够用其定义的属性，比如EType
			if err1.EType == 1 {
				fmt.Printf("Radius %0.2f is less than zero", err1.radius)
				return
			}

		}
		fmt.Println(err)
		return
	}
	fmt.Printf("Area of rectangle1 %0.2f", area)
}
