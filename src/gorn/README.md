https://github.com/roylee0704/gron

Cron Jobs in Go

分析总结：

1. 主要分两部分，一是schedule部分，二是cron具体实现Job。schedule是对时间的逻辑，实现Next，cron根据Next来定时执行
2.
2. 可以new cron实例：c， c.Start 的时候，就会新起一个goroutine，循环跑run方法，run方法中 第一步是轮训c.entries，重新配置 Next值，如果Next满足条件，则会 go entry.Job.Run()，新建goroutine非阻塞的执行具体Job逻辑
3. cron的多goroutine模型是简单的，一个主goroutine，不断派生子goroutine 执行job
4. 通过JobFunc func() 对外部传进来的func() 参数进行wrap adapter to Job，这个技巧很cool
5. schedule.go 的Every方法是比较重要的一个方法，用于定义 每隔多久执行操作，它得到的是每隔多少秒的时间
6. AtSchedule 内嵌了Schedule，所以，在传 Schedule的方法或类型的地方，也可以使用AtSchedule
