package main

import "testing"
import " pool"
import "pjob"

func BenchmarkConcurrent(b *testing.B) {
	collector := pool.StartDispatcher(WORKER_COUNT) // start up worker pool

	for n := 0; n < b.N; n++ {
		for i, job := range pjob.CreateJobs(20) {
			collector.Job <- pool.Job{Job: job, ID: i}
		}
	}
}

func BenchmarkNonconcurrent(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, job := range pjob.CreateJobs(20) {
			pjob.DoWork(job, 1)
		}
	}
}
