package job

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// create random string
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// create list of jobs
func CreateJobs(amount int) []string {
	var jobs []string

	for i := 0; i < amount; i++ {
		jobs = append(jobs, RandStringRunes(8))
	}
	return jobs
}

// mimics any type of job that can be run concurrently
func DoJob(word string, id int) {
	h := fnv.New32a()
	h.Write([]byte(word))
	time.Sleep(time.Second)
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("worker [%d] - created hash [%d] from word [%s]\n", id, h.Sum32(), word)
	}
}
