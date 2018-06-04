package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"strings"
	"sync"
	"time"

	fstream "github.com/adamdrake/gofstream"
)

func hash(s []byte) int {
	h := fnv.New64a()
	h.Write(s)
	return int(h.Sum64())
}

func worker(recs chan string, fields []string, w, n *[]float64, D, alpha float64, loss *float64, count *int, wg *sync.WaitGroup) {
	defer wg.Done()
	for r := range recs {
		request := strings.Split(r, ",")

		*count++
		y := 0
		if request[1] == "1" {
			y = 1
		}

		request = request[2:]            // ignore label and id
		x := make([]int, len(request)+1) // need length plus one for zero at front
		x[0] = 0
		for i, v := range request {
			hashResult := hash([]byte(fields[i]+v)) % int(D)
			x[i+1] = int(math.Abs(float64(hashResult)))
		}

		// Get the prediction for the given request (now transformed to hash values)
		wTx := 0.0
		for _, v := range x {
			wTx += (*w)[v]
		}
		p := 1.0 / (1.0 + math.Exp(-math.Max(math.Min(wTx, 20.0), -20.0)))

		// Update the loss
		p = math.Max(math.Min(p, 1.-math.Pow(10, -12)), math.Pow(10, -12))
		if y == 1 {
			*loss += -math.Log(p)
		} else {
			*loss += -math.Log(1.0 - p)
		}

		// Update the weights
		for _, v := range x {
			(*w)[v] = (*w)[v] - (p-float64(y))*alpha/(math.Sqrt((*n)[v])+1.0)
			(*n)[v]++
		}
	}
}

func main() {
	start := time.Now()
	D := math.Pow(2, 20) // number of weights use for learning
	alpha := 0.1         // learning rate for sgd optimization
	w := make([]float64, int(D))
	n := make([]float64, int(D))
	loss := 0.0
	count := 0

	data, _ := fstream.New("../train.csv", 10000)
	fields := strings.Split(<-data, ",")
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(data, fields, &w, &n, D, alpha, &loss, &count, &wg)
	}
	wg.Wait()
	fmt.Println("Run time is", time.Since(start))
	fmt.Println("loss", loss/float64(count))
	fmt.Println("RPS", float64(count)/time.Since(start).Seconds())
}
