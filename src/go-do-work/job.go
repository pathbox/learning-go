package gdw

type Job interface {
	DoWork()
}

type batchedJob struct {
	batched Job
	name    string
}
