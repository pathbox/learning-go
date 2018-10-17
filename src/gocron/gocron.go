package gocron

import (
	"errors"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var loc = time.Local

func ChangeLoc(newLocation *time.Location) {
	loc = newLocation
}

// Max number of jobs, hack it if you need
const MAXJOBNUM = 10000

type Job struct {
	// pause interval * unit bettween runs
	interval uint64

	// the job jobFunc to run, funcs[jobFunc]
	jobFunc string
	// time uints, ,e.g.  'minutes', 'hours'...
	uint string
	// optional time at which this job runs
	atTime string

	// datetime of last run
	lastRun time.Time
	// datetime of next run
	nextRun time.Time
	// cache the period between last an next run
	period time.Duration

	// Specific day of the week to start on
	startDay time.Weekday

	// Map for the function task store
	funcs map[string]interface{}

	// Map for function and params of function
	fparams map[string]([]interface{})
}

// Create a new job with the time interval
func NewJob(intervel uint64) *Job {
	return &Job{
		intervel,
		"", "", "",
		time.Unix(0, 0),
		time.Unix(0, 0), 0,
		time.Sunday,
		make(map[string]interface{}),
		make(map[string]([]interface{})),
	}
}

// True if job should be run now
func (j *Job) shouldRun() bool {
	return time.Now().After(j.nextRun) //  the time of next run is OK？ j.nextRun < time.Now()
}

// Run the job and immediately reschedule it
func (j *Job) run() (result []reflect.Value, err error) {
	f := reflect.ValueOf(j.funcs[j.jobFunc]) // 通过reflect，能够解析func属性，如：参数、函数名等
	params := j.fparams[j.jobFunc]
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of param is not adapted")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param) // 将参数reflect化后，放到in
	}
	result = f.Call(in) // 通过反射reflect的方式，函数，接收参数，然后执行
	j.lastRun = time.Now()
	j.scheduleNextRun()
	return
}

// for given function fn, get the name of function.
func getFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf((fn)).Pointer()).Name()
}

// Specifies the jobFunc that should be called every time the job runs

func (j *Job) Do(jobFun interface{}, params ...interface{}) {
	typ := reflect.TypeOf(jobFun)   // 反射化得到类型
	if typ.Kind() != reflect.Func { // 判断类型是否是 function
		panic("only function can be schedule into the job queue.")
	}

	fname := getFunctionName(jobFun)
	j.funcs[fname] = jobFun // func store to j.funcs
	j.fparams[fname] = params
	j.jobFunc = fname
	//schedule the next run
	j.scheduleNextRun()
}

func formatTime(t string) (hour, min int, err error) {
	var er = errors.Ner("time format error")
	ts := strings.Split(t, ":")
	if len(ts) != 2 {
		err = er
		return
	}

	hour, err = strconv.Atoi(ts[0])
	if err != nil {
		return
	}
	min, err = strconv.Atoi(ts[1])
	if err != nil {
		return
	}

	if hour < 0 || hour > 23 || min < 0 || min > 59 {
		err = er
		return
	}
	return hour, min, nil
}

//	s.Every(1).Day().At("10:30").Do(task)
//	s.Every(1).Monday().At("10:30").Do(task)

func (j *Job) At(t string) *Job {
	hour, min, err := formatTime(t)
	if err != nil {
		panic(err)
	}

	// time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	mock := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), int(hour), int(min), 0, 0, loc)

	if j.uint == "days" {
		if time.Now().After(mock) {
			j.lastRun = mock
		} else {
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, hour, min, 0, 0, loc)
		}
	} else if j.uint == "weeks" {
		if j.startDay != time.Now().Weekday() || (time.Now().After(mock) && j.startDay == time.Now().Weekday()) {
			i := mock.Weekday() - j.startDay
			if i < 0 {
				i = 7 + i
			}
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-int(i), hour, min, 0, 0, loc)
		} else {
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, hour, min, 0, 0, loc)
		}
	}
	return j
}

//Compute the instant when this job should run next
func (j *Job) scheduleNextRun() {
	if j.lastRun == time.Unix(0, 0) {
		if j.uint == "weeks" {
			i := time.Now().Weekday() - j.startDay
			if i < 0 {
				i = 7 + i
			}
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-int(i), 0, 0, 0, 0, loc)
		} else {
			j.lastRun = time.Now()
		}
	}
	if j.period != 0 {
		j.nextRun = j.lastRun.Add(j.period * time.Second)
	} else {
		switch j.unit {
		case "minutes":
			j.period = time.Duration(j.interval * 60)
			break
		case "hours":
			j.period = time.Duration(j.interval * 60 * 60)
			break
		case "days":
			j.period = time.Duration(j.interval * 60 * 60 * 24)
			break
		case "weeks":
			j.period = time.Duration(j.interval * 60 * 60 * 24 * 7)
			break
		case "seconds":
			j.period = time.Duration(j.interval)
		}
		j.nextRun = j.lastRun.Add(j.period * time.Second)
	}
}

// NextScheduledTime returns the time of when this job is to run next
func (j *Job) NextScheduledTime() time.Time {
	return j.nextRun
}

// the follow functions set the job's unit with seconds,minutes,hours...

// Set the unit with second
func (j *Job) Second() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Seconds()
	return
}

// Set the unit with seconds
func (j *Job) Seconds() (job *Job) {
	j.unit = "seconds"
	return j
}

// Set the unit  with minute, which interval is 1
func (j *Job) Minute() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Minutes()
	return
}

//set the unit with minute
func (j *Job) Minutes() (job *Job) {
	j.unit = "minutes"
	return j
}

//set the unit with hour, which interval is 1
func (j *Job) Hour() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Hours()
	return
}

// Set the unit with hours
func (j *Job) Hours() (job *Job) {
	j.unit = "hours"
	return j
}

// Set the job's unit with day, which interval is 1
func (j *Job) Day() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Days()
	return
}

// Set the job's unit with days
func (j *Job) Days() *Job {
	j.unit = "days"
	return j
}

// s.Every(1).Monday().Do(task)
// Set the start day with Monday
func (j *Job) Monday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 1
	job = j.Weeks()
	return
}

// Set the start day with Tuesday
func (j *Job) Tuesday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 2
	job = j.Weeks()
	return
}

// Set the start day woth Wednesday
func (j *Job) Wednesday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 3
	job = j.Weeks()
	return
}

// Set the start day with thursday
func (j *Job) Thursday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 4
	job = j.Weeks()
	return
}

// Set the start day with friday
func (j *Job) Friday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 5
	job = j.Weeks()
	return
}

// Set the start day with saturday
func (j *Job) Saturday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 6
	job = j.Weeks()
	return
}

// Set the start day with sunday
func (j *Job) Sunday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 0
	job = j.Weeks()
	return
}

//Set the units as weeks
func (j *Job) Weeks() *Job {
	j.unit = "weeks"
	return j
}

// Class Scheduler, the only data member is the list of jobs
// jobs 队列
type Scheduler struct {
	// Array store jobs
	jobs [MAXJOBNUM]*Job // job数组 在Scheduler层面，就是在调度[MAXJOBNUM]*Job 这个job数组了，在job层面就是对时间属性的操作和函数，参数的验证

	// Size of jobs which jobs holding.
	size int
}

// Scheduler implements the sort.Interface{} for sorting jobs, by the time nextRun

func (s *Scheduler) Len() int {
	return s.size
}

func (s *Scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

func (s *Scheduler) Less(i, j int) bool {
	return s.jobs[j].nextRun.After(s.jobs[i].nextRun)
}

// Create a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{[MAXJOBNUM]*Job{}, 0}
}

// Get the current runnalbe jobs, which shouldRun is True
func (s *Scheduler) getRunnableJobs() (running_jobs [MAXJOBNUM]*Job, n int) {
	runnableJobs := [MAXJOBNUM]*Job{}
	n = 0
	sort.Sort(s) // 每次对Scheduler的全局jobs进行进行时间的排序，时间小(说明时间到了)排在前面
	for i := 0; i < s.size; i++ {
		if s.jobs[i].shouldRun() { // 如果next运行时间已经小于当前时间，说明时间到了，可以运行了
			runnableJobs[n] = s.jobs[i]
			n++
		} else {
			break // 不满足时，则不用接下去比较了，因为之前已经排序了，后面的job都是还没到运行时间的
		}
	}
	return runnableJobs, n
}

// Datetime when the next job should run
func (s *Scheduler) NextRun() (*Job, time.Time) {
	if s.size <= 0 {
		return nil, time.Now()
	}
	sort.Sort(s)
	return s.jobs[0], s.jobs[0].nextRun
}

// Schedule a new periodic job，在这里将job存到 Scheduler的全局 jobs list中
func (s *Scheduler) Every(interval uint64) *Job {
	job := NewJob(interval)
	s.jobs[s.size] = job
	s.size++
	return job
}

// Run all the jobs that are scheduled to run
func (s *Scheduler) RunPending() {
	runnableJobs, n := s.getRunnableJobs()
	if n != 0 {
		for i := 0; i < n; i++ {
			runnableJobs[i].run()
		}
	}
}

// Run all jobs regardless if they are scheduled to run or not
func (s *Scheduler) RunAll() {
	for i := 0; i < s.size; i++ {
		s.jobs[i].run()
	}
}

// Run all jobs with delay seconds
func (s *Scheduler) RunAllwithDelay(d int) {
	for i := 0; i < s.size; i++ {
		s.jobs[i].run()
		time.Sleep(time.Duration(d))
	}
}

// Remove specific job j
func (s *Scheduler) Remove(j interface{}) {
	i := 0
	for ; i < s.size; i++ {
		if s.jobs[i].jobFunc == getFunctionname(j) {
			break
		}
	}

	for j := (i + 1); j < s.size; j++ {
		s.jobs[i] = s.jobs[j]
		i++
	}
	s.size = s.size - 1
}

// Delete all jobs
func (s *Scheduler) Clear() {
	for i := 0; i < s.size; i++ {
		s.jobs[i] = nil
	}
	s.size = 0
}

// Start all the pending jobs
// Add seconds ticker
func (s *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.RunPending() // 每秒会取获得可运行的job，然后进行运行
			case <-stopped:
				return
			}
		}
	}()

	return stopped
}

// The following methods are shortcuts for not having to
// create a Schduler instance

var defaultScheduler = NewScheduler()
var jobs = defaultScheduler.jobs

// Schedule a new periodic job
func Every(interval uint64) *Job {
	return defaultScheduler.Every(interval)
}

// Run all jobs that are scheduled to run
//
// Please note that it is *intended behavior that run_pending()
// does not run missed jobs*. For example, if you've registered a job
// that should run every minute and you only call run_pending()
// in one hour increments then your job won't be run 60 times in
// between but only once.
func RunPending() {
	defaultScheduler.RunPending()
}

// Run all jobs regardless if they are scheduled to run or not.
func RunAll() {
	defaultScheduler.RunAll()
}

// Run all the jobs with a delay in seconds
//
// A delay of `delay` seconds is added between each job. This can help
// to distribute the system load generated by the jobs more evenly over
// time.
func RunAllwithDelay(d int) {
	defaultScheduler.RunAllwithDelay(d)
}

// Run all jobs that are scheduled to run
func Start() chan bool {
	return defaultScheduler.Start()
}

// Clear
func Clear() {
	defaultScheduler.Clear()
}

// Remove
func Remove(j interface{}) {
	defaultScheduler.Remove(j)
}

// NextRun gets the next running time
func NextRun() (job *Job, time time.Time) {
	return defaultScheduler.NextRun()
}
