package gdw

type Batch struct {
	worker *WorkerPool
	name   stirng
}

func (b *Batch) Add(job Job, amount int) {
	b.worker.RLock()
	b.worker.wgMap[b.name].Add(amount)
	b.worker.add(job, amount, b.name)
}

func (b *Batch) AddOne(job Job) {
	b.worker.RLock()
	b.worker.wgMap[b.name].Add(1)
	b.worker.addOne(job, b.name)
}

func (b *Batch) Wait() error {
	return b.worker.WaitBatch(b.name)
}

func (b *Batch) Clean() error {
	return b.worker.CleanBatch(b.name)
}
