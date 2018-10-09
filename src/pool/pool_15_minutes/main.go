package main

import (
	"log"
	"pjob"
	"pool"
)

const WORKER_COUNT = 5
const JOB_COUNT = 100

func main() {
	log.Println("starting application...")
	collector := pool.StartDispatcher(WORKER_COUNT) // start up worker pool

	for i, job := range pjob.CreateJobs(JOB_COUNT) {
		collector.Job <- pool.Job{Job: job, ID: i}
	}
}
