package main

import "fmt"

func main() {
	var amap [20][20]string

	amap = [20][20]string{{"#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#"}, {}}
	i, j, sum, x, y, count, p, q := 0, 0, 0, 0, 0, 0, 0, 0
	for i = 0; i <= 12; i++ {
		for j = 0; j <= 12; j++ {
			if amap[i][j] == "." { // 判断这个点是不是平地，是平地才可以放炸弹。 一次起始点 i j不变，循环查找的是x y
				sum = 0 //可以消灭的敌人数
				// 向上统计可以消灭的敌人数
				x, y = i, j
				for amap[x][y] != "#" { // 判断不是墙，如果不是墙几继续
					if amap[x][y] == "G" { //如果当前点是敌人，则进行计数
						sum++
					}
					x-- // 继续向上统计
				}

				// 向下统计可以消灭的敌人数
				x, y = i, j
				for amap[x][y] != "#" { // 判断不是墙，如果不是墙几继续
					if amap[x][y] == "G" { //如果当前点是敌人，则进行计数
						sum++
					}
					x++ // 继续向下统计
				}

				// 向左统计可以消灭的敌人数
				x, y = i, j
				for amap[x][y] != "#" { // 判断不是墙，如果不是墙几继续
					if amap[x][y] == "G" { //如果当前点是敌人，则进行计数
						sum++
					}
					y-- // 继续向左统计
				}

				// 向右统计可以消灭的敌人数
				x, y = i, j
				for amap[x][y] != "#" { // 判断不是墙，如果不是墙几继续
					if amap[x][y] == "G" { //如果当前点是敌人，则进行计数
						sum++
					}
					y++ // 继续向右统计
				}

				// 更新count的值
				if sum > count {
					count = sum
					p = i
					q = j
				}
			}
		}
	}

	fmt.Println(p, q, count) // 在 p q 点，可以消灭做多的敌人count个
}
