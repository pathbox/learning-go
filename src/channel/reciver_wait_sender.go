// reciver in main wait for the sender from other goroutinue
package main

import (
	"fmt"
)

type UpdateOp struct {
	key int
}

func applyUpdate(data map[int]string, op UpdateOp) {
	data[op.key] = op.value
}

func main() {
	m := make(map[int]string)
	m[2] = "asdf"

	ch := make(chan UpdateOp)

	go updateOperation(ch)

	applyUpdate(m, <-ch) // reciver 操作会一直阻塞等待 sender的值
	// go updateOperation(ch chan UpdateOp) 方法不能在applyUpdate之后
	fmt.Println(m[2])

}

func updateOperation(ch chan UpdateOp) {
	ch <- UpdateOp{2, "New Value"}
}

// channel和mutex的本质区别是一件事情是在"被调用处"做还是在"调用处"做，channel称不上比mutex好。
// 解释：
// channel代表的是两点间的关系，而很多现实问题是多点的，这个时候使用channel最自然的解决方案就是：有一个角色负责操作某件事情或某个资源，其他线程都通过channel向这个角色发号施令。如果我们在程序中设置N个角色，让它们各司其职，那么程序就能分类有序地运转下去。所以使用channel的潜台词就是把程序划分为不同的角色。但mutex不同，每个线程想干什么事时，直接去获取干这件事的权力(I mean mutex, obviously)，然后直接把事干了。
// channel固然直观，但是有代价：额外的上下文切换。做成任何事情都得等到被调用处被调度，处理，回复，调用处才能继续。这个再怎么优化，再怎么尊重cache locality，相比无竞争时的锁也没什么卵用。而好的并发代码的无竞争比例都挺高的，所以这是channel必须吃的硬亏。另外一个现实是：用channel的代码比用mutex更难写。由于业务一致性的限制，一些资源往往被绑定在一起，所以一个角色很可能身兼数职，但它做一件事情时便无法做另一件事情，而事情又有优先级。各种打断、跳出、继续形成的最终代码异常复杂。当然这些问题也可以映射到mutex的世界中，但在使用channel时这种倾向更明显。
// 所以说到这份上，是否使用channel更多是考虑性能和口味。性能不那么重要，大家又觉得使用channel好理解，那就用。如果数据是单向的，那么也会用(buffered) channel，这时channel的作用更多是队列。
