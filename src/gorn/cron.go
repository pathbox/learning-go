package gorn

import (
	"sort"
	"time"
)

// Entry consists of a schedule and the job to be executed on that schedule.
type Entry struct {
	Schedule Schedule
	Job      Job

	// the next time the job will run. This is zero time if Cron has not been
	// started or invalid schedule.
	Next time.Time

	// the last time the job was run. This is zero time if the job has not been
	// run.
	Prev time.Time
}

// byTime is a handy wrapper to chronologically sort entries.
type byTime []*Entry

func (b byTime) Len() int      { return len(b) }
func (b byTime) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// Less reports `earliest` time i should sort before j.
// zero time is not `earliest` time.
func (b byTime) Less(i, j int) bool {
	if b[i].Next.IsZero() {
		return false
	}
	if b[j].Next.IsZero() {
		return true
	}

	return b[i].Next.Before(b[j].Next)
}

// Job is the interface that wraps the basic Run method.
//
// Run executes the underlying func.
type Job interface { // Job是一个接口
	Run()
}

// Cron provides a convenient interface for scheduling job such as to clean-up
// database entry every month.
//
// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may also be started, stopped and the entries
// may be inspected
type Cron struct {
	entries []*Entry
	running bool
	add     chan *Entry
	stop    chan struct{}
}

func New() *Cron {
	return &Cron{
		stop: make(chan struct{}),
		add:  make(chan *Entry),
	}
}

// start 表示 开启一个goroutine 执行run循环处理
// 每一个 cron实例，执行run方法，占用一个goroutine
func (c *Cron) Start() {
	c.running = ture
	go c.run()
}

// Add appends schedule, job to entries.
//
// if cron instant is not running, adding to entries is trivial
func (c *Cron) Add(s Schedule, j Job) {
	entry := &Entry{
		Schedule: s,
		Job:      j,
	}

	if !c.running { //如果 cron没有在执行
		c.entries = append(c.entries, entry)
		return
	}

	c.add <- entry // 如果cron在执行，则通过add chan 接收 新的entry
}

// AddFunc registers the Job function for the given Schedule
// 为 cron 注册 job func， job func就是你需要具体实现的调度逻辑，然后作为参数传入
func (c *Cron) AddFunc(s Schedule, j func()) {
	c.Add(s, JobFunc(j)) // JobFUnc() is cool
}

func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.running = false
	c.stop <- struct{}{}
}

var after = time.After

// run the scheduler...
//
// It needs to be private as it's responsible of synchronizing a critical
// shared state: `running`.
func (c *Cron) run() {
	var effective time.Time
	now := time.Now().Local()

	// to figure next trig time for entries, referenced fron now
	for _, e := range c.entries { // 轮训 cron中的所有entry
		e.Next = e.Schedule.Next(now) // 很具 entry 各自的schedule配置，重新设置Next时间值
	}

	for {
		sort.Sort(byTime(c.entries))
		if len(c.entries) > 0 {
			effective = c.entries[0].Next
		} else {
			effective = now.AddDate(15, 0, 0) // to prevent phantom jobs.
		}

		select {
		case now = <-after(effective.Sub(now)):
			for _, entry := range c.entries {
				if entry.Next != effective {
					break
				}
				entry.Prev = now
				entry.Next = entry.Schedule.Next(now)
				go entry.Job.Run() // 每次新起一个goroutine异步执行 具体的job.Run
			}
		case e := <-c.add: //通过 c.add将entry加入到c.entries,避免在运行的时候，写操作c.entries 造成竞态问题
			e.Next = e.Schedule.Next(time.Now())
			c.entries = append(c.entries, e)
		case <-c.stop:
			return // 结束 cron 的run goroutine
		}
	}
}

// Entries returns cron etn
func (c Cron) Entries() []*Entry {
	return c.entries
}

// JobFunc is an adapter to allow the use of ordinary functions as gron.Job
// If f is a function with the appropriate signature, JobFunc(f) is a handler
// that calls f.
//
// todo: possibly func with params? maybe not needed.
// JobFunc 这个 adapter 的设计很有技巧。首先，JobFunc也是一种类型，和struct是同等的，类型就是 func()
// JobFunc 实现了Run 方法，复合Job接口，这样就实现了Job接口，所以可以赋值给Corn
// 定义一个 j := func(){...具体代码}，j是一个func类型的变量，要将j转为 JobFunc类型，只要 jNow := JobFunc(j)就OK啦！ 这个和 i := 32; i64 := int64(i) 原理是一样的。 因为j的底层类型是fun()和 JobFunc的底层类型func()是一样的，Golang是允许他们可以互相转换的。方法就是 Type() 这样，非常简便。反过来说，如果两个变量的底层类型是不能互相转换的，就无法使用上述的方法
type JobFunc func()

// Run calls j()
func (j JobFunc) Run() {
	j()
}
